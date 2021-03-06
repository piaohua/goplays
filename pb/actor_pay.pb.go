// Code generated by protoc-gen-gogo. DO NOT EDIT.
// source: actor_pay.proto

package pb

import proto "github.com/gogo/protobuf/proto"
import fmt "fmt"
import math "math"

import strings "strings"
import reflect "reflect"

import io "io"

// Reference imports to suppress errors if they are not otherwise used.
var _ = proto.Marshal
var _ = fmt.Errorf
var _ = math.Inf

// 苹果支付
type ApplePay struct {
	Trade string `protobuf:"bytes,1,opt,name=Trade,proto3" json:"Trade,omitempty"`
}

func (m *ApplePay) Reset()                    { *m = ApplePay{} }
func (*ApplePay) ProtoMessage()               {}
func (*ApplePay) Descriptor() ([]byte, []int) { return fileDescriptorActorPay, []int{0} }

func (m *ApplePay) GetTrade() string {
	if m != nil {
		return m.Trade
	}
	return ""
}

type ApplePaid struct {
	Result bool `protobuf:"varint,1,opt,name=Result,proto3" json:"Result,omitempty"`
}

func (m *ApplePaid) Reset()                    { *m = ApplePaid{} }
func (*ApplePaid) ProtoMessage()               {}
func (*ApplePaid) Descriptor() ([]byte, []int) { return fileDescriptorActorPay, []int{1} }

func (m *ApplePaid) GetResult() bool {
	if m != nil {
		return m.Result
	}
	return false
}

// 微信支付主动回调或主动查询发货
type WxpayCallback struct {
	Result string `protobuf:"bytes,1,opt,name=Result,proto3" json:"Result,omitempty"`
}

func (m *WxpayCallback) Reset()                    { *m = WxpayCallback{} }
func (*WxpayCallback) ProtoMessage()               {}
func (*WxpayCallback) Descriptor() ([]byte, []int) { return fileDescriptorActorPay, []int{2} }

func (m *WxpayCallback) GetResult() string {
	if m != nil {
		return m.Result
	}
	return ""
}

// 发货
type WxpayGoods struct {
	Userid  string `protobuf:"bytes,1,opt,name=Userid,proto3" json:"Userid,omitempty"`
	Orderid string `protobuf:"bytes,2,opt,name=Orderid,proto3" json:"Orderid,omitempty"`
	Money   uint32 `protobuf:"varint,3,opt,name=Money,proto3" json:"Money,omitempty"`
	Diamond uint32 `protobuf:"varint,4,opt,name=Diamond,proto3" json:"Diamond,omitempty"`
	First   int32  `protobuf:"varint,5,opt,name=First,proto3" json:"First,omitempty"`
}

func (m *WxpayGoods) Reset()                    { *m = WxpayGoods{} }
func (*WxpayGoods) ProtoMessage()               {}
func (*WxpayGoods) Descriptor() ([]byte, []int) { return fileDescriptorActorPay, []int{3} }

func (m *WxpayGoods) GetUserid() string {
	if m != nil {
		return m.Userid
	}
	return ""
}

func (m *WxpayGoods) GetOrderid() string {
	if m != nil {
		return m.Orderid
	}
	return ""
}

func (m *WxpayGoods) GetMoney() uint32 {
	if m != nil {
		return m.Money
	}
	return 0
}

func (m *WxpayGoods) GetDiamond() uint32 {
	if m != nil {
		return m.Diamond
	}
	return 0
}

func (m *WxpayGoods) GetFirst() int32 {
	if m != nil {
		return m.First
	}
	return 0
}

func init() {
	proto.RegisterType((*ApplePay)(nil), "pb.ApplePay")
	proto.RegisterType((*ApplePaid)(nil), "pb.ApplePaid")
	proto.RegisterType((*WxpayCallback)(nil), "pb.WxpayCallback")
	proto.RegisterType((*WxpayGoods)(nil), "pb.WxpayGoods")
}
func (this *ApplePay) Equal(that interface{}) bool {
	if that == nil {
		return this == nil
	}

	that1, ok := that.(*ApplePay)
	if !ok {
		that2, ok := that.(ApplePay)
		if ok {
			that1 = &that2
		} else {
			return false
		}
	}
	if that1 == nil {
		return this == nil
	} else if this == nil {
		return false
	}
	if this.Trade != that1.Trade {
		return false
	}
	return true
}
func (this *ApplePaid) Equal(that interface{}) bool {
	if that == nil {
		return this == nil
	}

	that1, ok := that.(*ApplePaid)
	if !ok {
		that2, ok := that.(ApplePaid)
		if ok {
			that1 = &that2
		} else {
			return false
		}
	}
	if that1 == nil {
		return this == nil
	} else if this == nil {
		return false
	}
	if this.Result != that1.Result {
		return false
	}
	return true
}
func (this *WxpayCallback) Equal(that interface{}) bool {
	if that == nil {
		return this == nil
	}

	that1, ok := that.(*WxpayCallback)
	if !ok {
		that2, ok := that.(WxpayCallback)
		if ok {
			that1 = &that2
		} else {
			return false
		}
	}
	if that1 == nil {
		return this == nil
	} else if this == nil {
		return false
	}
	if this.Result != that1.Result {
		return false
	}
	return true
}
func (this *WxpayGoods) Equal(that interface{}) bool {
	if that == nil {
		return this == nil
	}

	that1, ok := that.(*WxpayGoods)
	if !ok {
		that2, ok := that.(WxpayGoods)
		if ok {
			that1 = &that2
		} else {
			return false
		}
	}
	if that1 == nil {
		return this == nil
	} else if this == nil {
		return false
	}
	if this.Userid != that1.Userid {
		return false
	}
	if this.Orderid != that1.Orderid {
		return false
	}
	if this.Money != that1.Money {
		return false
	}
	if this.Diamond != that1.Diamond {
		return false
	}
	if this.First != that1.First {
		return false
	}
	return true
}
func (this *ApplePay) GoString() string {
	if this == nil {
		return "nil"
	}
	s := make([]string, 0, 5)
	s = append(s, "&pb.ApplePay{")
	s = append(s, "Trade: "+fmt.Sprintf("%#v", this.Trade)+",\n")
	s = append(s, "}")
	return strings.Join(s, "")
}
func (this *ApplePaid) GoString() string {
	if this == nil {
		return "nil"
	}
	s := make([]string, 0, 5)
	s = append(s, "&pb.ApplePaid{")
	s = append(s, "Result: "+fmt.Sprintf("%#v", this.Result)+",\n")
	s = append(s, "}")
	return strings.Join(s, "")
}
func (this *WxpayCallback) GoString() string {
	if this == nil {
		return "nil"
	}
	s := make([]string, 0, 5)
	s = append(s, "&pb.WxpayCallback{")
	s = append(s, "Result: "+fmt.Sprintf("%#v", this.Result)+",\n")
	s = append(s, "}")
	return strings.Join(s, "")
}
func (this *WxpayGoods) GoString() string {
	if this == nil {
		return "nil"
	}
	s := make([]string, 0, 9)
	s = append(s, "&pb.WxpayGoods{")
	s = append(s, "Userid: "+fmt.Sprintf("%#v", this.Userid)+",\n")
	s = append(s, "Orderid: "+fmt.Sprintf("%#v", this.Orderid)+",\n")
	s = append(s, "Money: "+fmt.Sprintf("%#v", this.Money)+",\n")
	s = append(s, "Diamond: "+fmt.Sprintf("%#v", this.Diamond)+",\n")
	s = append(s, "First: "+fmt.Sprintf("%#v", this.First)+",\n")
	s = append(s, "}")
	return strings.Join(s, "")
}
func valueToGoStringActorPay(v interface{}, typ string) string {
	rv := reflect.ValueOf(v)
	if rv.IsNil() {
		return "nil"
	}
	pv := reflect.Indirect(rv).Interface()
	return fmt.Sprintf("func(v %v) *%v { return &v } ( %#v )", typ, typ, pv)
}
func (m *ApplePay) Marshal() (dAtA []byte, err error) {
	size := m.Size()
	dAtA = make([]byte, size)
	n, err := m.MarshalTo(dAtA)
	if err != nil {
		return nil, err
	}
	return dAtA[:n], nil
}

func (m *ApplePay) MarshalTo(dAtA []byte) (int, error) {
	var i int
	_ = i
	var l int
	_ = l
	if len(m.Trade) > 0 {
		dAtA[i] = 0xa
		i++
		i = encodeVarintActorPay(dAtA, i, uint64(len(m.Trade)))
		i += copy(dAtA[i:], m.Trade)
	}
	return i, nil
}

func (m *ApplePaid) Marshal() (dAtA []byte, err error) {
	size := m.Size()
	dAtA = make([]byte, size)
	n, err := m.MarshalTo(dAtA)
	if err != nil {
		return nil, err
	}
	return dAtA[:n], nil
}

func (m *ApplePaid) MarshalTo(dAtA []byte) (int, error) {
	var i int
	_ = i
	var l int
	_ = l
	if m.Result {
		dAtA[i] = 0x8
		i++
		if m.Result {
			dAtA[i] = 1
		} else {
			dAtA[i] = 0
		}
		i++
	}
	return i, nil
}

func (m *WxpayCallback) Marshal() (dAtA []byte, err error) {
	size := m.Size()
	dAtA = make([]byte, size)
	n, err := m.MarshalTo(dAtA)
	if err != nil {
		return nil, err
	}
	return dAtA[:n], nil
}

func (m *WxpayCallback) MarshalTo(dAtA []byte) (int, error) {
	var i int
	_ = i
	var l int
	_ = l
	if len(m.Result) > 0 {
		dAtA[i] = 0xa
		i++
		i = encodeVarintActorPay(dAtA, i, uint64(len(m.Result)))
		i += copy(dAtA[i:], m.Result)
	}
	return i, nil
}

func (m *WxpayGoods) Marshal() (dAtA []byte, err error) {
	size := m.Size()
	dAtA = make([]byte, size)
	n, err := m.MarshalTo(dAtA)
	if err != nil {
		return nil, err
	}
	return dAtA[:n], nil
}

func (m *WxpayGoods) MarshalTo(dAtA []byte) (int, error) {
	var i int
	_ = i
	var l int
	_ = l
	if len(m.Userid) > 0 {
		dAtA[i] = 0xa
		i++
		i = encodeVarintActorPay(dAtA, i, uint64(len(m.Userid)))
		i += copy(dAtA[i:], m.Userid)
	}
	if len(m.Orderid) > 0 {
		dAtA[i] = 0x12
		i++
		i = encodeVarintActorPay(dAtA, i, uint64(len(m.Orderid)))
		i += copy(dAtA[i:], m.Orderid)
	}
	if m.Money != 0 {
		dAtA[i] = 0x18
		i++
		i = encodeVarintActorPay(dAtA, i, uint64(m.Money))
	}
	if m.Diamond != 0 {
		dAtA[i] = 0x20
		i++
		i = encodeVarintActorPay(dAtA, i, uint64(m.Diamond))
	}
	if m.First != 0 {
		dAtA[i] = 0x28
		i++
		i = encodeVarintActorPay(dAtA, i, uint64(m.First))
	}
	return i, nil
}

func encodeVarintActorPay(dAtA []byte, offset int, v uint64) int {
	for v >= 1<<7 {
		dAtA[offset] = uint8(v&0x7f | 0x80)
		v >>= 7
		offset++
	}
	dAtA[offset] = uint8(v)
	return offset + 1
}
func (m *ApplePay) Size() (n int) {
	var l int
	_ = l
	l = len(m.Trade)
	if l > 0 {
		n += 1 + l + sovActorPay(uint64(l))
	}
	return n
}

func (m *ApplePaid) Size() (n int) {
	var l int
	_ = l
	if m.Result {
		n += 2
	}
	return n
}

func (m *WxpayCallback) Size() (n int) {
	var l int
	_ = l
	l = len(m.Result)
	if l > 0 {
		n += 1 + l + sovActorPay(uint64(l))
	}
	return n
}

func (m *WxpayGoods) Size() (n int) {
	var l int
	_ = l
	l = len(m.Userid)
	if l > 0 {
		n += 1 + l + sovActorPay(uint64(l))
	}
	l = len(m.Orderid)
	if l > 0 {
		n += 1 + l + sovActorPay(uint64(l))
	}
	if m.Money != 0 {
		n += 1 + sovActorPay(uint64(m.Money))
	}
	if m.Diamond != 0 {
		n += 1 + sovActorPay(uint64(m.Diamond))
	}
	if m.First != 0 {
		n += 1 + sovActorPay(uint64(m.First))
	}
	return n
}

func sovActorPay(x uint64) (n int) {
	for {
		n++
		x >>= 7
		if x == 0 {
			break
		}
	}
	return n
}
func sozActorPay(x uint64) (n int) {
	return sovActorPay(uint64((x << 1) ^ uint64((int64(x) >> 63))))
}
func (this *ApplePay) String() string {
	if this == nil {
		return "nil"
	}
	s := strings.Join([]string{`&ApplePay{`,
		`Trade:` + fmt.Sprintf("%v", this.Trade) + `,`,
		`}`,
	}, "")
	return s
}
func (this *ApplePaid) String() string {
	if this == nil {
		return "nil"
	}
	s := strings.Join([]string{`&ApplePaid{`,
		`Result:` + fmt.Sprintf("%v", this.Result) + `,`,
		`}`,
	}, "")
	return s
}
func (this *WxpayCallback) String() string {
	if this == nil {
		return "nil"
	}
	s := strings.Join([]string{`&WxpayCallback{`,
		`Result:` + fmt.Sprintf("%v", this.Result) + `,`,
		`}`,
	}, "")
	return s
}
func (this *WxpayGoods) String() string {
	if this == nil {
		return "nil"
	}
	s := strings.Join([]string{`&WxpayGoods{`,
		`Userid:` + fmt.Sprintf("%v", this.Userid) + `,`,
		`Orderid:` + fmt.Sprintf("%v", this.Orderid) + `,`,
		`Money:` + fmt.Sprintf("%v", this.Money) + `,`,
		`Diamond:` + fmt.Sprintf("%v", this.Diamond) + `,`,
		`First:` + fmt.Sprintf("%v", this.First) + `,`,
		`}`,
	}, "")
	return s
}
func valueToStringActorPay(v interface{}) string {
	rv := reflect.ValueOf(v)
	if rv.IsNil() {
		return "nil"
	}
	pv := reflect.Indirect(rv).Interface()
	return fmt.Sprintf("*%v", pv)
}
func (m *ApplePay) Unmarshal(dAtA []byte) error {
	l := len(dAtA)
	iNdEx := 0
	for iNdEx < l {
		preIndex := iNdEx
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return ErrIntOverflowActorPay
			}
			if iNdEx >= l {
				return io.ErrUnexpectedEOF
			}
			b := dAtA[iNdEx]
			iNdEx++
			wire |= (uint64(b) & 0x7F) << shift
			if b < 0x80 {
				break
			}
		}
		fieldNum := int32(wire >> 3)
		wireType := int(wire & 0x7)
		if wireType == 4 {
			return fmt.Errorf("proto: ApplePay: wiretype end group for non-group")
		}
		if fieldNum <= 0 {
			return fmt.Errorf("proto: ApplePay: illegal tag %d (wire type %d)", fieldNum, wire)
		}
		switch fieldNum {
		case 1:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Trade", wireType)
			}
			var stringLen uint64
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowActorPay
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				stringLen |= (uint64(b) & 0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			intStringLen := int(stringLen)
			if intStringLen < 0 {
				return ErrInvalidLengthActorPay
			}
			postIndex := iNdEx + intStringLen
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.Trade = string(dAtA[iNdEx:postIndex])
			iNdEx = postIndex
		default:
			iNdEx = preIndex
			skippy, err := skipActorPay(dAtA[iNdEx:])
			if err != nil {
				return err
			}
			if skippy < 0 {
				return ErrInvalidLengthActorPay
			}
			if (iNdEx + skippy) > l {
				return io.ErrUnexpectedEOF
			}
			iNdEx += skippy
		}
	}

	if iNdEx > l {
		return io.ErrUnexpectedEOF
	}
	return nil
}
func (m *ApplePaid) Unmarshal(dAtA []byte) error {
	l := len(dAtA)
	iNdEx := 0
	for iNdEx < l {
		preIndex := iNdEx
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return ErrIntOverflowActorPay
			}
			if iNdEx >= l {
				return io.ErrUnexpectedEOF
			}
			b := dAtA[iNdEx]
			iNdEx++
			wire |= (uint64(b) & 0x7F) << shift
			if b < 0x80 {
				break
			}
		}
		fieldNum := int32(wire >> 3)
		wireType := int(wire & 0x7)
		if wireType == 4 {
			return fmt.Errorf("proto: ApplePaid: wiretype end group for non-group")
		}
		if fieldNum <= 0 {
			return fmt.Errorf("proto: ApplePaid: illegal tag %d (wire type %d)", fieldNum, wire)
		}
		switch fieldNum {
		case 1:
			if wireType != 0 {
				return fmt.Errorf("proto: wrong wireType = %d for field Result", wireType)
			}
			var v int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowActorPay
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				v |= (int(b) & 0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			m.Result = bool(v != 0)
		default:
			iNdEx = preIndex
			skippy, err := skipActorPay(dAtA[iNdEx:])
			if err != nil {
				return err
			}
			if skippy < 0 {
				return ErrInvalidLengthActorPay
			}
			if (iNdEx + skippy) > l {
				return io.ErrUnexpectedEOF
			}
			iNdEx += skippy
		}
	}

	if iNdEx > l {
		return io.ErrUnexpectedEOF
	}
	return nil
}
func (m *WxpayCallback) Unmarshal(dAtA []byte) error {
	l := len(dAtA)
	iNdEx := 0
	for iNdEx < l {
		preIndex := iNdEx
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return ErrIntOverflowActorPay
			}
			if iNdEx >= l {
				return io.ErrUnexpectedEOF
			}
			b := dAtA[iNdEx]
			iNdEx++
			wire |= (uint64(b) & 0x7F) << shift
			if b < 0x80 {
				break
			}
		}
		fieldNum := int32(wire >> 3)
		wireType := int(wire & 0x7)
		if wireType == 4 {
			return fmt.Errorf("proto: WxpayCallback: wiretype end group for non-group")
		}
		if fieldNum <= 0 {
			return fmt.Errorf("proto: WxpayCallback: illegal tag %d (wire type %d)", fieldNum, wire)
		}
		switch fieldNum {
		case 1:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Result", wireType)
			}
			var stringLen uint64
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowActorPay
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				stringLen |= (uint64(b) & 0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			intStringLen := int(stringLen)
			if intStringLen < 0 {
				return ErrInvalidLengthActorPay
			}
			postIndex := iNdEx + intStringLen
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.Result = string(dAtA[iNdEx:postIndex])
			iNdEx = postIndex
		default:
			iNdEx = preIndex
			skippy, err := skipActorPay(dAtA[iNdEx:])
			if err != nil {
				return err
			}
			if skippy < 0 {
				return ErrInvalidLengthActorPay
			}
			if (iNdEx + skippy) > l {
				return io.ErrUnexpectedEOF
			}
			iNdEx += skippy
		}
	}

	if iNdEx > l {
		return io.ErrUnexpectedEOF
	}
	return nil
}
func (m *WxpayGoods) Unmarshal(dAtA []byte) error {
	l := len(dAtA)
	iNdEx := 0
	for iNdEx < l {
		preIndex := iNdEx
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return ErrIntOverflowActorPay
			}
			if iNdEx >= l {
				return io.ErrUnexpectedEOF
			}
			b := dAtA[iNdEx]
			iNdEx++
			wire |= (uint64(b) & 0x7F) << shift
			if b < 0x80 {
				break
			}
		}
		fieldNum := int32(wire >> 3)
		wireType := int(wire & 0x7)
		if wireType == 4 {
			return fmt.Errorf("proto: WxpayGoods: wiretype end group for non-group")
		}
		if fieldNum <= 0 {
			return fmt.Errorf("proto: WxpayGoods: illegal tag %d (wire type %d)", fieldNum, wire)
		}
		switch fieldNum {
		case 1:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Userid", wireType)
			}
			var stringLen uint64
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowActorPay
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				stringLen |= (uint64(b) & 0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			intStringLen := int(stringLen)
			if intStringLen < 0 {
				return ErrInvalidLengthActorPay
			}
			postIndex := iNdEx + intStringLen
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.Userid = string(dAtA[iNdEx:postIndex])
			iNdEx = postIndex
		case 2:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Orderid", wireType)
			}
			var stringLen uint64
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowActorPay
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				stringLen |= (uint64(b) & 0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			intStringLen := int(stringLen)
			if intStringLen < 0 {
				return ErrInvalidLengthActorPay
			}
			postIndex := iNdEx + intStringLen
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.Orderid = string(dAtA[iNdEx:postIndex])
			iNdEx = postIndex
		case 3:
			if wireType != 0 {
				return fmt.Errorf("proto: wrong wireType = %d for field Money", wireType)
			}
			m.Money = 0
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowActorPay
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				m.Money |= (uint32(b) & 0x7F) << shift
				if b < 0x80 {
					break
				}
			}
		case 4:
			if wireType != 0 {
				return fmt.Errorf("proto: wrong wireType = %d for field Diamond", wireType)
			}
			m.Diamond = 0
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowActorPay
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				m.Diamond |= (uint32(b) & 0x7F) << shift
				if b < 0x80 {
					break
				}
			}
		case 5:
			if wireType != 0 {
				return fmt.Errorf("proto: wrong wireType = %d for field First", wireType)
			}
			m.First = 0
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowActorPay
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				m.First |= (int32(b) & 0x7F) << shift
				if b < 0x80 {
					break
				}
			}
		default:
			iNdEx = preIndex
			skippy, err := skipActorPay(dAtA[iNdEx:])
			if err != nil {
				return err
			}
			if skippy < 0 {
				return ErrInvalidLengthActorPay
			}
			if (iNdEx + skippy) > l {
				return io.ErrUnexpectedEOF
			}
			iNdEx += skippy
		}
	}

	if iNdEx > l {
		return io.ErrUnexpectedEOF
	}
	return nil
}
func skipActorPay(dAtA []byte) (n int, err error) {
	l := len(dAtA)
	iNdEx := 0
	for iNdEx < l {
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return 0, ErrIntOverflowActorPay
			}
			if iNdEx >= l {
				return 0, io.ErrUnexpectedEOF
			}
			b := dAtA[iNdEx]
			iNdEx++
			wire |= (uint64(b) & 0x7F) << shift
			if b < 0x80 {
				break
			}
		}
		wireType := int(wire & 0x7)
		switch wireType {
		case 0:
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return 0, ErrIntOverflowActorPay
				}
				if iNdEx >= l {
					return 0, io.ErrUnexpectedEOF
				}
				iNdEx++
				if dAtA[iNdEx-1] < 0x80 {
					break
				}
			}
			return iNdEx, nil
		case 1:
			iNdEx += 8
			return iNdEx, nil
		case 2:
			var length int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return 0, ErrIntOverflowActorPay
				}
				if iNdEx >= l {
					return 0, io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				length |= (int(b) & 0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			iNdEx += length
			if length < 0 {
				return 0, ErrInvalidLengthActorPay
			}
			return iNdEx, nil
		case 3:
			for {
				var innerWire uint64
				var start int = iNdEx
				for shift := uint(0); ; shift += 7 {
					if shift >= 64 {
						return 0, ErrIntOverflowActorPay
					}
					if iNdEx >= l {
						return 0, io.ErrUnexpectedEOF
					}
					b := dAtA[iNdEx]
					iNdEx++
					innerWire |= (uint64(b) & 0x7F) << shift
					if b < 0x80 {
						break
					}
				}
				innerWireType := int(innerWire & 0x7)
				if innerWireType == 4 {
					break
				}
				next, err := skipActorPay(dAtA[start:])
				if err != nil {
					return 0, err
				}
				iNdEx = start + next
			}
			return iNdEx, nil
		case 4:
			return iNdEx, nil
		case 5:
			iNdEx += 4
			return iNdEx, nil
		default:
			return 0, fmt.Errorf("proto: illegal wireType %d", wireType)
		}
	}
	panic("unreachable")
}

var (
	ErrInvalidLengthActorPay = fmt.Errorf("proto: negative length found during unmarshaling")
	ErrIntOverflowActorPay   = fmt.Errorf("proto: integer overflow")
)

func init() { proto.RegisterFile("actor_pay.proto", fileDescriptorActorPay) }

var fileDescriptorActorPay = []byte{
	// 257 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0xe2, 0xe2, 0x4f, 0x4c, 0x2e, 0xc9,
	0x2f, 0x8a, 0x2f, 0x48, 0xac, 0xd4, 0x2b, 0x28, 0xca, 0x2f, 0xc9, 0x17, 0x62, 0x2a, 0x48, 0x52,
	0x52, 0xe0, 0xe2, 0x70, 0x2c, 0x28, 0xc8, 0x49, 0x0d, 0x48, 0xac, 0x14, 0x12, 0xe1, 0x62, 0x0d,
	0x29, 0x4a, 0x4c, 0x49, 0x95, 0x60, 0x54, 0x60, 0xd4, 0xe0, 0x0c, 0x82, 0x70, 0x94, 0x94, 0xb9,
	0x38, 0xa1, 0x2a, 0x32, 0x53, 0x84, 0xc4, 0xb8, 0xd8, 0x82, 0x52, 0x8b, 0x4b, 0x73, 0x4a, 0xc0,
	0x6a, 0x38, 0x82, 0xa0, 0x3c, 0x25, 0x75, 0x2e, 0xde, 0xf0, 0x8a, 0x82, 0xc4, 0x4a, 0xe7, 0xc4,
	0x9c, 0x9c, 0xa4, 0xc4, 0xe4, 0x6c, 0x34, 0x85, 0x9c, 0x70, 0x85, 0x2d, 0x8c, 0x5c, 0x5c, 0x60,
	0x95, 0xee, 0xf9, 0xf9, 0x29, 0xc5, 0x20, 0x65, 0xa1, 0xc5, 0xa9, 0x45, 0x99, 0x29, 0x30, 0x65,
	0x10, 0x9e, 0x90, 0x04, 0x17, 0xbb, 0x7f, 0x51, 0x0a, 0x58, 0x82, 0x09, 0x2c, 0x01, 0xe3, 0x82,
	0x1c, 0xe9, 0x9b, 0x9f, 0x97, 0x5a, 0x29, 0xc1, 0xac, 0xc0, 0xa8, 0xc1, 0x1b, 0x04, 0xe1, 0x80,
	0xd4, 0xbb, 0x64, 0x26, 0xe6, 0xe6, 0xe7, 0xa5, 0x48, 0xb0, 0x80, 0xc5, 0x61, 0x5c, 0x90, 0x7a,
	0xb7, 0xcc, 0xa2, 0xe2, 0x12, 0x09, 0x56, 0x05, 0x46, 0x0d, 0xd6, 0x20, 0x08, 0xc7, 0x49, 0xe7,
	0xc2, 0x43, 0x39, 0x86, 0x1b, 0x0f, 0xe5, 0x18, 0x3e, 0x3c, 0x94, 0x63, 0x6c, 0x78, 0x24, 0xc7,
	0xb8, 0xe2, 0x91, 0x1c, 0xe3, 0x89, 0x47, 0x72, 0x8c, 0x17, 0x1e, 0xc9, 0x31, 0x3e, 0x78, 0x24,
	0xc7, 0xf8, 0xe2, 0x91, 0x1c, 0xc3, 0x87, 0x47, 0x72, 0x8c, 0x13, 0x1e, 0xcb, 0x31, 0x24, 0xb1,
	0x81, 0xc3, 0xcb, 0x18, 0x10, 0x00, 0x00, 0xff, 0xff, 0x24, 0x78, 0x15, 0x09, 0x42, 0x01, 0x00,
	0x00,
}
