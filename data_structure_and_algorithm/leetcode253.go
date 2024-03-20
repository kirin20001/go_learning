package data_structure_and_algorithm

import "sort"

func minMeetingRooms(intervals [][]int) int {
	timeline := make([]int, 2*len(intervals))
	for _, v := range intervals {
		timeline = append(timeline, v[0]*10+1)
		timeline = append(timeline, v[1]*10)
	}
	sort.Ints(timeline)
	var numRoom, curMeeting int
	for _, v := range timeline {
		if v%10 == 0 {
			curMeeting--
		}
		if v %10 == 1 {
			curMeeting++
			numRoom = max(numRoom, curMeeting)
		}
	}

	return numRoom
}

