package shouldwater

import "testing"


func TestNotEnoughData(t *testing.T) {
	_, err := ShouldWater([]WeatherRecord{}, []WeatherRecord{})

	if err == nil {
		t.Errorf("Should return err when there is not enough data")
	}
}
