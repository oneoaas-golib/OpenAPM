// AUTOGENERATED FILE: easyjson marshaller/unmarshallers.

package service

import (
	json "encoding/json"
	jlexer "github.com/mailru/easyjson/jlexer"
	jwriter "github.com/mailru/easyjson/jwriter"
)

// suppress unused package warning
var (
	_ = json.RawMessage{}
	_ = jlexer.Lexer{}
	_ = jwriter.Writer{}
)

func easyjson_5c87105d_decode_github_com_aiyun_vgo_vgo_alarm_service_AlertData(in *jlexer.Lexer, out *AlertData) {
	if in.IsNull() {
		in.Skip()
		return
	}
	in.Delim('{')
	for !in.IsDelim('}') {
		key := in.UnsafeString()
		in.WantColon()
		if in.IsNull() {
			in.Skip()
			in.WantComma()
			continue
		}
		switch key {
		case "id":
			out.ID = string(in.String())
		case "gid":
			out.GroupID = string(in.String())
		case "v":
			out.Value = float64(in.Float64())
		default:
			in.SkipRecursive()
		}
		in.WantComma()
	}
	in.Delim('}')
}
func easyjson_5c87105d_encode_github_com_aiyun_vgo_vgo_alarm_service_AlertData(out *jwriter.Writer, in AlertData) {
	out.RawByte('{')
	first := true
	_ = first
	if !first {
		out.RawByte(',')
	}
	first = false
	out.RawString("\"id\":")
	out.String(string(in.ID))
	if !first {
		out.RawByte(',')
	}
	first = false
	out.RawString("\"gid\":")
	out.String(string(in.GroupID))
	if !first {
		out.RawByte(',')
	}
	first = false
	out.RawString("\"v\":")
	out.Float64(float64(in.Value))
	out.RawByte('}')
}
func (v AlertData) MarshalJSON() ([]byte, error) {
	w := jwriter.Writer{}
	easyjson_5c87105d_encode_github_com_aiyun_vgo_vgo_alarm_service_AlertData(&w, v)
	return w.Buffer.BuildBytes(), w.Error
}
func (v AlertData) MarshalEasyJSON(w *jwriter.Writer) {
	easyjson_5c87105d_encode_github_com_aiyun_vgo_vgo_alarm_service_AlertData(w, v)
}
func (v *AlertData) UnmarshalJSON(data []byte) error {
	r := jlexer.Lexer{Data: data}
	easyjson_5c87105d_decode_github_com_aiyun_vgo_vgo_alarm_service_AlertData(&r, v)
	return r.Error()
}
func (v *AlertData) UnmarshalEasyJSON(l *jlexer.Lexer) {
	easyjson_5c87105d_decode_github_com_aiyun_vgo_vgo_alarm_service_AlertData(l, v)
}
