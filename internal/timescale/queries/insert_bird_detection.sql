INSERT INTO sensors.birdnet (
    time, source_node, source, begin_time, end_time,
    species_code, scientific_name, common_name, confidence,
    latitude, longitude, threshold, sensitivity
) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13)