package engine

import (
	"fmt"
	"testing"
)

func TestUnpackPackFunctions(t *testing.T) {
	expectedMove := Move(MOVE_MASK)
	expectedEval := MAX_INT
	expectedDepth := int8(127)
	expectedType := LowerBound
	expectedAge := uint16(1023)
	data := Pack(expectedMove, expectedEval, expectedDepth, expectedType, expectedAge)
	actualMove, actualEval, actualDepth, actualType, actualAge := Unpack(data)

	if actualMove != expectedMove {
		t.Errorf("Unexpected move: %s", fmt.Sprintf("Expected: %d, Got %d\n", expectedMove, actualMove))
	}

	if actualEval != expectedEval {
		t.Errorf("Unexpected eval: %s", fmt.Sprintf("Expected: %d, Got %d\n", expectedEval, actualEval))
	}

	if actualDepth != expectedDepth {
		t.Errorf("Unexpected depth: %s", fmt.Sprintf("Expected: %d, Got %d\n", expectedDepth, actualDepth))
	}

	if actualType != expectedType {
		t.Errorf("Unexpected type: %s", fmt.Sprintf("Expected: %d, Got %d\n", expectedType, actualType))
	}

	if actualAge != expectedAge {
		t.Errorf("Unexpected age: %s", fmt.Sprintf("Expected: %d, Got %d\n", expectedAge, actualAge))
	}

	expectedMove = NewMove(E1, E2, WhiteKing, BlackKing, Queen, KingSideCastle|QueenSideCastle|Capture|EnPassant)
	expectedEval = int16(0)
	expectedDepth = int8(50)
	expectedType = Exact
	expectedAge = uint16(512)
	data = Pack(expectedMove, expectedEval, expectedDepth, expectedType, expectedAge)
	actualMove, actualEval, actualDepth, actualType, actualAge = Unpack(data)

	if actualMove != expectedMove {
		t.Errorf("Unexpected move: %s", fmt.Sprintf("Expected: %d, Got %d\n", expectedMove, actualMove))
	}

	if actualEval != expectedEval {
		t.Errorf("Unexpected eval: %s", fmt.Sprintf("Expected: %d, Got %d\n", expectedEval, actualEval))
	}

	if actualDepth != expectedDepth {
		t.Errorf("Unexpected depth: %s", fmt.Sprintf("Expected: %d, Got %d\n", expectedDepth, actualDepth))
	}

	if actualType != expectedType {
		t.Errorf("Unexpected type: %s", fmt.Sprintf("Expected: %d, Got %d\n", expectedType, actualType))
	}

	if actualAge != expectedAge {
		t.Errorf("Unexpected age: %s", fmt.Sprintf("Expected: %d, Got %d\n", expectedAge, actualAge))
	}
}
