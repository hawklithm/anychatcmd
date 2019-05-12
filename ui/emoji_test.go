package ui

import (
	"fmt"
	"testing"
)

func Test_GetUUID(t *testing.T) {
	fmt.Println(ConvertToEmoji(
		"12124[微笑] sdfkjasdjflaj[哈哈12]1e0dfajsdlsdjfklsjf"))
	fmt.Println(TranslateEmoji(
		"124215235<span class=\"emoji emoji1f61d\"></span" +
			">24edflaksjdlkgasjdlkgjladg" +
		"<span class=\"emoji emoji1f633\"></span>dfasdfasdf"))
}
