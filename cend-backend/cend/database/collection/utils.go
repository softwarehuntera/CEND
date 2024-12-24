package collection

import (
	"log/slog"
	"runtime"
	"fmt"
)


func LogInfo(message string) {
	pc, file, line, ok := runtime.Caller(1)
	if !ok {
		slog.Info("Could not retrieve caller information", "message", message)
		return
	}
	funcName := runtime.FuncForPC(pc).Name()
	slog.Info(message, "file", file, "line", line, "function", funcName)
}

func Equal(actual *Collection, expected *Collection) bool {
	// Compare each key in the lookup table
	for key, expectedValue := range *expected.lookupTable {
		actualValue, exists := (*actual.lookupTable)[key]
		if !exists {
			LogInfo(fmt.Sprintf("Expected key %v in lookup table but was not found", key))
			return false
		}

		// Check frequency
		if actualValue.count != expectedValue.count {
			LogInfo(fmt.Sprintf("Frequency mismatch for key %v: expected %v, got %v", key, expectedValue.count, actualValue.count))
			return false
		}

		// Check document IDs
		if len(actualValue.docIDs) != len(expectedValue.docIDs) {
			LogInfo(fmt.Sprintf("DocIDs mismatch for key %v: expected %v, got %v", key, expectedValue.docIDs, actualValue.docIDs))
			return false
		}
		for docID := range expectedValue.docIDs {
			if _, exists := actualValue.docIDs[docID]; !exists {
				LogInfo(fmt.Sprintf("Expected docID %v for key %v not found in actual docIDs", docID, key))
				return false
			}
		}
	}

	// Compare documents
	if expected.documents.Length() != actual.documents.Length() {
		LogInfo(fmt.Sprintf("Documents mismatch. actual=%v, expected=%v", *expected.documents,*actual.documents ))
		return false
	}
	
	for docID, expectedDoc := range expected.documents.Documents() {
		actualDoc := actual.documents.Get(docID)
		if actualDoc == nil {
			LogInfo(fmt.Sprintf("Expected docID %v not found in actual documents", docID))
			return false
		}
		actualDocStr := actualDoc.String()
		if actualDocStr != expectedDoc.String() {
			LogInfo(fmt.Sprintf("Mismatch in documents: expected docID %v to have content %v, got %v", docID, expectedDoc, actualDoc))
			return false
		}
	}
	return true
}
