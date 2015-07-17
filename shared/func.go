/* This file is part of GHTS.

GHTS is free software: you can redistribute it and/or modify
it under the terms of the GNU Lesser General Public License as published by
the Free Software Foundation, either version 3 of the License, or
(at your option) any later version.

GHTS is distributed in the hope that it will be useful,
but WITHOUT ANY WARRANTY; without even the implied warranty of
MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
GNU Lesser General Public License for more details.

You should have received a copy of the GNU Lesser General Public License
along with GHTS.  If not, see <http://www.gnu.org/licenses/>.

@author: UnHa Kim <unha.kim.ghts@gmail.com> */

package shared

import (
	"reflect"
	"runtime"
)

func F문자열_모음2인터페이스_모음(문자열_모음 []string) []interface{} {
	if 문자열_모음 == nil {
		return nil
	}

	인터페이스_모음 := make([]interface{}, len(문자열_모음))

	for i := 0; i < len(문자열_모음); i++ {
		인터페이스_모음[i] = 문자열_모음[i]
	}

	return 인터페이스_모음
}

func F인터페이스_모음2문자열_모음(인터페이스_모음 []interface{}) []string {
	if 인터페이스_모음 == nil {
		return nil
	}

	문자열_모음 := make([]string, len(인터페이스_모음))

	for i := 0; i < len(인터페이스_모음); i++ {
		문자열_모음[i] = F포맷된_문자열("%v", 인터페이스_모음[i])
	}

	return 문자열_모음
}

func F타입_이름(i interface{}) string {
	return reflect.TypeOf(i).Name()
}

func F문자열_복사(문자열 string) string {
	return (문자열 + " ")[:len(문자열)]
}

// 이하 최대 스레드 수량 관련 함수

func F단일_스레드_모드() { runtime.GOMAXPROCS(1) }
func F멀티_스레드_모드() { runtime.GOMAXPROCS(runtime.NumCPU()) }

func F단일_스레드_모드임() bool {
	if runtime.GOMAXPROCS(-1) == 1 {
		return true
	} else {
		return false
	}
}

func F멀티_스레드_모드임() bool { return !F단일_스레드_모드임() }

// 이하 종료 시 존재하는 모든 Go루틴 정리(혹은 종료) 관련 함수 모음
var ch공통_종료_채널 = make(chan S비어있는_구조체)

func F공통_종료_채널() chan S비어있는_구조체 {
	return ch공통_종료_채널
}

func F등록된_Go루틴_종료() {
	close(ch공통_종료_채널)
}

func F_nil에러() error { return nil }
