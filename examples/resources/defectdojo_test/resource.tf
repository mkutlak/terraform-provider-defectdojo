resource "defectdojo_test" "example" {
  test_type    = 1
  engagement   = 1
  target_start = "2025-01-01T00:00:00Z"
  target_end   = "2025-01-31T23:59:59Z"
  title        = "SAST Scan"
}
