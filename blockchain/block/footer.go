// Copyright (c) 2019 IoTeX Foundation
// This source code is provided 'as is' and no warranties are given as to title or non-infringement, merchantability
// or fitness for purpose and, to the extent permitted by law, all liability for your use of the code is disclaimed.
// This source code is governed by Apache License 2.0 that can be found in the LICENSE file.

package block

import (
	"time"

	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/types/known/timestamppb"

	"github.com/iotexproject/iotex-core/v2/endorsement"
	"github.com/iotexproject/iotex-proto/golang/iotextypes"
)

// Footer defines a set of proof of this block
type Footer struct {
	endorsements []*endorsement.Endorsement
	commitTime   time.Time
}

// Proto converts BlockFooter
func (f *Footer) Proto() *iotextypes.BlockFooter {
	pb := iotextypes.BlockFooter{}
	commitTime := timestamppb.New(f.commitTime)
	pb.Timestamp = commitTime
	pb.Endorsements = []*iotextypes.Endorsement{}
	for _, en := range f.endorsements {
		pb.Endorsements = append(pb.Endorsements, en.Proto())
	}
	return &pb
}

// ConvertFromBlockFooterPb converts BlockFooter to BlockFooter
func (f *Footer) ConvertFromBlockFooterPb(pb *iotextypes.BlockFooter) error {
	if pb == nil {
		return nil
	}
	if err := pb.GetTimestamp().CheckValid(); err != nil {
		return err
	}
	commitTime := pb.GetTimestamp().AsTime()
	f.commitTime = commitTime
	pbEndorsements := pb.GetEndorsements()
	if pbEndorsements == nil {
		return nil
	}
	f.endorsements = []*endorsement.Endorsement{}
	for _, ePb := range pbEndorsements {
		e := &endorsement.Endorsement{}
		if err := e.LoadProto(ePb); err != nil {
			return err
		}
		f.endorsements = append(f.endorsements, e)
	}

	return nil
}

// CommitTime returns the timestamp the block was committed
func (f *Footer) CommitTime() time.Time {
	return f.commitTime
}

// Endorsements returns the number of commit endorsements froms delegates
func (f *Footer) Endorsements() []*endorsement.Endorsement {
	return f.endorsements
}

// Serialize returns the serialized byte stream of the block footer
func (f *Footer) Serialize() ([]byte, error) {
	return proto.Marshal(f.Proto())
}

// Deserialize loads from the serialized byte stream
func (f *Footer) Deserialize(buf []byte) error {
	pb := &iotextypes.BlockFooter{}
	if err := proto.Unmarshal(buf, pb); err != nil {
		return err
	}
	return f.ConvertFromBlockFooterPb(pb)
}
