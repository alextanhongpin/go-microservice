package model

// NOTE: Just use *time.Time in a struct...
// type NullableTime struct {
//         time.Time `json:",omitempty"`
// }
//
// func (n NullableTime) MarshalJSON() ([]byte, error) {
//         if n.Time.IsZero() {
//                 // Optional: return []byte("null"), but it will not be omitted
//                 return []byte(`""`), nil
//         }
//
//         return n.Time.MarshalJSON()
// }
//
// // Scan must be a pointer.
// func (n *NullableTime) Scan(value interface{}) error {
//         if value == nil {
//                 return nil
//         }
//         var ok bool
//         n.Time, ok = value.(time.Time)
//         if !ok {
//                 return errors.New("failed to scan model.NullableTime")
//         }
//         return nil
// }
//
// func (n NullableTime) Value() (driver.Value, error) {
//         return n.Time.Format(time.RFC3339Nano), nil
// }
