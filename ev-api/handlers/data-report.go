package handlers

import (
	"context"
	"ev-api/config"
	"math"
	"strconv"
	"time"

	"github.com/gofiber/fiber/v2"
)

func DataReport(c *fiber.Ctx) error {
	pageNo, err := strconv.Atoi(c.Query("pageNo", "1"))
	if err != nil || pageNo < 1 {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid pageNo parameter",
		})
	}
	pageSize, err := strconv.Atoi(c.Query("pageSize", "10"))
	if err != nil || pageSize < 1 {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid pageSize parameter",
		})
	}
	offset := (pageNo - 1) * pageSize

	finalQuery := `
		WITH aggregated AS (
			SELECT
				td.id::text AS id,
				(
					SELECT jsonb_object_agg(key, 
						CASE 
							WHEN value ~ '^-?\d+(\.\d+)?$' THEN to_jsonb(value::numeric)
							ELSE to_jsonb(value)
						END
					) 
					FROM jsonb_each_text(td.gps_data)
				) AS gps_data,
				TO_CHAR(TO_TIMESTAMP('20' || (td.gps_data->>'event_time'), 'YYYYMMDDHH24MISS'), 'DD-MM-YYYY"T"HH24:MI:SS') AS event_time,
				jsonb_object_agg(
					CASE 
						WHEN array_length(regexp_split_to_array(trim(both '"' FROM signal.value::text), ' '), 1) > 1 THEN
							signal.key || ' (' ||
								array_to_string(
									(regexp_split_to_array(trim(both '"' FROM signal.value::text), ' '))[2:], 
									' '
								)
								|| ')'
						ELSE
							signal.key
					END,
					(split_part(trim(both '"' FROM signal.value::text), ' ', 1))::numeric
				) AS can_data
			FROM tbl_devicedata td
			CROSS JOIN LATERAL jsonb_each(td.can_data->'data') AS data
			CROSS JOIN LATERAL jsonb_each(data.value) AS signal
			GROUP BY td.id, td.gps_data
		)
		SELECT
			aggregated.*,
			COUNT(*) OVER() AS total_count
		FROM aggregated
		ORDER BY event_time ASC
		LIMIT $1 OFFSET $2
	`

	rows, err := config.DB.Query(context.Background(), finalQuery, pageSize, offset)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to execute query: " + err.Error(),
		})
	}
	defer rows.Close()

	colDescs := rows.FieldDescriptions()
	columns := make([]string, len(colDescs))
	for i, col := range colDescs {
		columns[i] = string(col.Name)
	}

	var (
		results    []map[string]interface{}
		totalCount int64
	)

	for rows.Next() {
		values := make([]interface{}, len(columns))
		valuePtrs := make([]interface{}, len(columns))
		for i := range values {
			valuePtrs[i] = &values[i]
		}

		if err := rows.Scan(valuePtrs...); err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": "Failed to scan row: " + err.Error(),
			})
		}

		rowData := make(map[string]interface{})
		for i, col := range columns {
			val := values[i]
			switch v := val.(type) {
			case time.Time:
				rowData[col] = v.Format(time.RFC3339)
			case []byte:
				rowData[col] = string(v)
			default:
				rowData[col] = v
			}

			if col == "total_count" && val != nil {
				totalCount = val.(int64)
			}
		}

		delete(rowData, "total_count")

		results = append(results, rowData)
	}

	totalPages := int(math.Ceil(float64(totalCount) / float64(pageSize)))

	return c.JSON(fiber.Map{
		"data":       	results,
		"current_page": pageNo,
		"page_size":  	pageSize,
		"total_pages":  totalPages,
	})
}
