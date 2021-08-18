// Code generated by easyjson for marshaling/unmarshaling. DO NOT EDIT.

package serde

import (
	json "encoding/json"
	easyjson "github.com/mailru/easyjson"
	jlexer "github.com/mailru/easyjson/jlexer"
	jwriter "github.com/mailru/easyjson/jwriter"
)

// suppress unused package warning
var (
	_ *json.RawMessage
	_ *jlexer.Lexer
	_ *jwriter.Writer
	_ easyjson.Marshaler
)

func easyjsonEeca4a30DecodeGithubComRidebeamGolangSkeletonPkgSerde(in *jlexer.Lexer, out *ExampleOutput) {
	isTopLevel := in.IsStart()
	if in.IsNull() {
		if isTopLevel {
			in.Consumed()
		}
		in.Skip()
		return
	}
	in.Delim('{')
	for !in.IsDelim('}') {
		key := in.UnsafeFieldName(false)
		in.WantColon()
		if in.IsNull() {
			in.Skip()
			in.WantComma()
			continue
		}
		switch key {
		case "ID":
			out.ID = string(in.String())
		case "foo_bar":
			if in.IsNull() {
				in.Skip()
			} else {
				in.Delim('{')
				if !in.IsDelim('}') {
					out.FooBar = make(map[string]string)
				} else {
					out.FooBar = nil
				}
				for !in.IsDelim('}') {
					key := string(in.String())
					in.WantColon()
					var v1 string
					v1 = string(in.String())
					(out.FooBar)[key] = v1
					in.WantComma()
				}
				in.Delim('}')
			}
		default:
			in.SkipRecursive()
		}
		in.WantComma()
	}
	in.Delim('}')
	if isTopLevel {
		in.Consumed()
	}
}
func easyjsonEeca4a30EncodeGithubComRidebeamGolangSkeletonPkgSerde(out *jwriter.Writer, in ExampleOutput) {
	out.RawByte('{')
	first := true
	_ = first
	if in.ID != "" {
		const prefix string = ",\"ID\":"
		first = false
		out.RawString(prefix[1:])
		out.String(string(in.ID))
	}
	if len(in.FooBar) != 0 {
		const prefix string = ",\"foo_bar\":"
		if first {
			first = false
			out.RawString(prefix[1:])
		} else {
			out.RawString(prefix)
		}
		{
			out.RawByte('{')
			v2First := true
			for v2Name, v2Value := range in.FooBar {
				if v2First {
					v2First = false
				} else {
					out.RawByte(',')
				}
				out.String(string(v2Name))
				out.RawByte(':')
				out.String(string(v2Value))
			}
			out.RawByte('}')
		}
	}
	out.RawByte('}')
}

// MarshalJSON supports json.Marshaler interface
func (v ExampleOutput) MarshalJSON() ([]byte, error) {
	w := jwriter.Writer{}
	easyjsonEeca4a30EncodeGithubComRidebeamGolangSkeletonPkgSerde(&w, v)
	return w.Buffer.BuildBytes(), w.Error
}

// MarshalEasyJSON supports easyjson.Marshaler interface
func (v ExampleOutput) MarshalEasyJSON(w *jwriter.Writer) {
	easyjsonEeca4a30EncodeGithubComRidebeamGolangSkeletonPkgSerde(w, v)
}

// UnmarshalJSON supports json.Unmarshaler interface
func (v *ExampleOutput) UnmarshalJSON(data []byte) error {
	r := jlexer.Lexer{Data: data}
	easyjsonEeca4a30DecodeGithubComRidebeamGolangSkeletonPkgSerde(&r, v)
	return r.Error()
}

// UnmarshalEasyJSON supports easyjson.Unmarshaler interface
func (v *ExampleOutput) UnmarshalEasyJSON(l *jlexer.Lexer) {
	easyjsonEeca4a30DecodeGithubComRidebeamGolangSkeletonPkgSerde(l, v)
}
func easyjsonEeca4a30DecodeGithubComRidebeamGolangSkeletonPkgSerde1(in *jlexer.Lexer, out *ExampleModel) {
	isTopLevel := in.IsStart()
	if in.IsNull() {
		if isTopLevel {
			in.Consumed()
		}
		in.Skip()
		return
	}
	in.Delim('{')
	for !in.IsDelim('}') {
		key := in.UnsafeFieldName(false)
		in.WantColon()
		if in.IsNull() {
			in.Skip()
			in.WantComma()
			continue
		}
		switch key {
		case "ID":
			out.ID = string(in.String())
		case "foo":
			out.Foo = string(in.String())
		case "bar":
			if in.IsNull() {
				in.Skip()
			} else {
				in.Delim('{')
				if !in.IsDelim('}') {
					out.Bar = make(map[string]string)
				} else {
					out.Bar = nil
				}
				for !in.IsDelim('}') {
					key := string(in.String())
					in.WantColon()
					var v3 string
					v3 = string(in.String())
					(out.Bar)[key] = v3
					in.WantComma()
				}
				in.Delim('}')
			}
		default:
			in.SkipRecursive()
		}
		in.WantComma()
	}
	in.Delim('}')
	if isTopLevel {
		in.Consumed()
	}
}
func easyjsonEeca4a30EncodeGithubComRidebeamGolangSkeletonPkgSerde1(out *jwriter.Writer, in ExampleModel) {
	out.RawByte('{')
	first := true
	_ = first
	if in.ID != "" {
		const prefix string = ",\"ID\":"
		first = false
		out.RawString(prefix[1:])
		out.String(string(in.ID))
	}
	if in.Foo != "" {
		const prefix string = ",\"foo\":"
		if first {
			first = false
			out.RawString(prefix[1:])
		} else {
			out.RawString(prefix)
		}
		out.String(string(in.Foo))
	}
	if len(in.Bar) != 0 {
		const prefix string = ",\"bar\":"
		if first {
			first = false
			out.RawString(prefix[1:])
		} else {
			out.RawString(prefix)
		}
		{
			out.RawByte('{')
			v4First := true
			for v4Name, v4Value := range in.Bar {
				if v4First {
					v4First = false
				} else {
					out.RawByte(',')
				}
				out.String(string(v4Name))
				out.RawByte(':')
				out.String(string(v4Value))
			}
			out.RawByte('}')
		}
	}
	out.RawByte('}')
}

// MarshalJSON supports json.Marshaler interface
func (v ExampleModel) MarshalJSON() ([]byte, error) {
	w := jwriter.Writer{}
	easyjsonEeca4a30EncodeGithubComRidebeamGolangSkeletonPkgSerde1(&w, v)
	return w.Buffer.BuildBytes(), w.Error
}

// MarshalEasyJSON supports easyjson.Marshaler interface
func (v ExampleModel) MarshalEasyJSON(w *jwriter.Writer) {
	easyjsonEeca4a30EncodeGithubComRidebeamGolangSkeletonPkgSerde1(w, v)
}

// UnmarshalJSON supports json.Unmarshaler interface
func (v *ExampleModel) UnmarshalJSON(data []byte) error {
	r := jlexer.Lexer{Data: data}
	easyjsonEeca4a30DecodeGithubComRidebeamGolangSkeletonPkgSerde1(&r, v)
	return r.Error()
}

// UnmarshalEasyJSON supports easyjson.Unmarshaler interface
func (v *ExampleModel) UnmarshalEasyJSON(l *jlexer.Lexer) {
	easyjsonEeca4a30DecodeGithubComRidebeamGolangSkeletonPkgSerde1(l, v)
}