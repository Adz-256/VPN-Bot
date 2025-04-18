package utils

import (
	"errors"
	"fmt"
	"strconv"
	"strings"
)

var (
	ErrorNoMask           = errors.New("no ip mask")
	ErrorMalformedIP      = errors.New("malformed ip")
	ErrorIPSubNetOverflow = errors.New("ip subnet overflow")
)

func IncrIP(ip string) (string, error) {
	if !strings.Contains(ip, "/") {
		return "", fmt.Errorf("%w: %s", ErrorNoMask, ip)
	}

	if !strings.Contains(ip, ".") {
		return "", fmt.Errorf("%w: %s", ErrorMalformedIP, ip)
	}

	ips := strings.Split(ip, ".")
	if len(ips) != 4 {
		return "", fmt.Errorf("%w: %s", ErrorMalformedIP, ip)
	}
	//получаем маску
	maskWithoctet := strings.Split(ips[len(ips)-1], "/")
	if len(maskWithoctet) != 2 {
		return "", fmt.Errorf("%w: %s", ErrorMalformedIP, ip)
	}
	mask := maskWithoctet[1]
	//убираем маску из последнего элемента ip
	ips[len(ips)-1] = ips[len(ips)-1][0 : len(ips[len(ips)-1])-len(mask)-1]

	for i := len(ips) - 1; i >= 0; i-- {
		ipInt, err := strconv.Atoi(ips[i])
		if err != nil {
			return "", fmt.Errorf("%w: %s", ErrorMalformedIP, ip)
		}

		ipInt++

		maskInt, err := strconv.Atoi(mask)
		if err != nil {
			return "", fmt.Errorf("%w: %s", ErrorMalformedIP, ip)
		}

		if i == (32-maskInt)/8 && ipInt > 256-2<<(maskInt%8)/2 {
			return "", ErrorIPSubNetOverflow
		}

		if ipInt >= 255 {
			ips[i] = "0"
			continue
		}

		ips[i] = strconv.Itoa(ipInt)
		break
	}

	return strings.Join(ips, ".") + "/" + mask, nil
}
