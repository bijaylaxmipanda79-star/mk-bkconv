package mihon

import (
	"compress/gzip"
	"io"
	"os"

	pb "github.com/galpt/mk-bkconv/proto/mihon"
	"google.golang.org/protobuf/proto"
)

// LoadBackup reads a Mihon backup file (.tachibk) using protoc-generated types.
// Auto-detects gzip compression and unmarshals using google.golang.org/protobuf.
func LoadBackup(path string) (*pb.Backup, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	// Read first two bytes to detect gzip (0x1f8b magic)
	hdr := make([]byte, 2)
	if _, err := f.Read(hdr); err != nil {
		return nil, err
	}
	if _, err := f.Seek(0, io.SeekStart); err != nil {
		return nil, err
	}

	var data []byte
	// Check for gzip magic bytes
	if hdr[0] == 0x1f && hdr[1] == 0x8b {
		gr, err := gzip.NewReader(f)
		if err != nil {
			return nil, err
		}
		defer gr.Close()
		data, err = io.ReadAll(gr)
		if err != nil {
			return nil, err
		}
	} else {
		data, err = io.ReadAll(f)
		if err != nil {
			return nil, err
		}
	}

	// Unmarshal using generated protobuf code
	backup := &pb.Backup{}
	if err := proto.Unmarshal(data, backup); err != nil {
		return nil, err
	}

	return backup, nil
}

// WriteBackup writes a Mihon backup using protoc-generated types.
// Marshals to protobuf and gzips the output.
func WriteBackup(path string, backup *pb.Backup) error {
	// Marshal using generated protobuf code
	data, err := proto.Marshal(backup)
	if err != nil {
		return err
	}

	// Create output file
	outf, err := os.Create(path)
	if err != nil {
		return err
	}
	defer outf.Close()

	// Gzip compress
	gw := gzip.NewWriter(outf)
	defer gw.Close()

	_, err = gw.Write(data)
	return err
}
