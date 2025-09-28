package main

import (
	"errors"
	"fmt"
)

var ErrInvalidModulo11 error = errors.New("invalid nhs number")

// modulo11 generates a modulo 11 addendum to a number based on the
// specification of NHS numbers described at Wikipedia.
//
// https://en.wikipedia.org/wiki/NHS_number
//
// The checksum is calculated by multiplying each of the first nine
// digits by 11 minus its position. Using the number 943 476 5919 as an
// example:
//   - The first digit is 9. This is multiplied by 10.
//   - The second digit is 4. This is multiplied by 9.
//   - And so on until the ninth digit (1) is multiplied by 2.
//   - The result of this calculation is summed. In this example:
//     (9×10) + (4×9) + (3×8) + (4×7) + (7×6) + (6×5) + (5×4) + (9×3) + (1×2) = 299.
//   - The remainder when dividing this number by 11 is calculated,
//     yielding a number in the range 0–10, which would be 2 in this case.
//   - Finally, this number is subtracted from 11 to give the checksum in
//     the range 1–11, in this case 9, which becomes the last digit of the
//     NHS number.
//   - A checksum of 11 is represented by 0 in the final NHS number. If
//     the checksum is 10 then the number is not valid.
func modulo11(number int) (int, error) {
	if number < 100000000 {
		return -1, fmt.Errorf("number below 100000000, got %d", number)
	}
	if number > 999999999 {
		return -1, fmt.Errorf("number above 999999999, got %d", number)
	}

	tempNo := number
	sum := 0

	for pos := 8; pos >= 0; pos-- {
		lastDigit := tempNo % 10
		tempNo /= 10 // remove last number
		weight := 10 - pos
		sum += lastDigit * weight
	}
	remainder := 11 - (sum % 11)
	switch remainder {
	case 10:
		return -1, ErrInvalidModulo11
	case 11:
		remainder = 0
	}
	finalNo := number*10 + remainder
	return finalNo, nil
}
