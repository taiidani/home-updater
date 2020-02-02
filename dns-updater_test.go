package main

import "testing"

func Test_extractDomain(t *testing.T) {
	tests := []struct {
		name     string
		domain   string
		wantApex string
		wantSub  string
	}{
		{
			name:     "Solo Apex",
			domain:   "apex.com",
			wantApex: "apex.com",
			wantSub:  "",
		},
		{
			name:     "Subdomain",
			domain:   "sub.apex.com",
			wantApex: "apex.com",
			wantSub:  "sub",
		},
		{
			name:     "Sub subdomain",
			domain:   "sub.sub.apex.com",
			wantApex: "apex.com",
			wantSub:  "sub.sub",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotApex, gotSub := extractDomain(tt.domain)
			if gotApex != tt.wantApex {
				t.Errorf("extractDomain() gotApex = %v, want %v", gotApex, tt.wantApex)
			}
			if gotSub != tt.wantSub {
				t.Errorf("extractDomain() gotSub = %v, want %v", gotSub, tt.wantSub)
			}
		})
	}
}
