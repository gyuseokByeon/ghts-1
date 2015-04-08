/*
This file is part of GHTS.

GHTS is free software: you can redistribute it and/or modify
it under the terms of the GNU General Public License as published by
the Free Software Foundation, either version 3 of the License, or
(at your option) any later version.

GHTS is distributed in the hope that it will be useful,
but WITHOUT ANY WARRANTY; without even the implied warranty of
MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
GNU General Public License for more details.

You should have received a copy of the GNU General Public License
along with GHTS.  If not, see <http://www.gnu.org/licenses/>.

Created on 2015. 4. 5.

@author: UnHa Kim <unha.kim@gh-system.com>
*/

package price_data_publish

import (
	공용 "github.com/gh-system/ghts/shared"
	zmq "github.com/pebbe/zmq4"
	"strconv"
	"testing"
	"time"
)

func TestF가격정보_배포_모듈_파이썬(테스트 *testing.T) {
    //체크포인트 := 1
    //공용.F체크포인트(&체크포인트, "테스트 시작")
    
    테스트_결과_수신_소켓, 에러 := zmq.NewSocket(zmq.REP)
	defer 테스트_결과_수신_소켓.Close()

	if 에러 != nil {
		공용.F문자열_출력("테스트_결과_수신_소켓 초기화 중 에러 발생. %s", 에러.Error())
		테스트.Fail()
	}
	
	테스트_결과_수신_소켓.Bind(공용.P테스트_결과_회신_주소)
	
	//공용.F체크포인트(&체크포인트, "테스트_결과_수신_소켓 초기화 완료")
    
	가격정보_배포횟수 := 1000
	구독_모듈_수량 := 10
	
	go F가격정보_배포_모듈()
	//공용.F체크포인트(&체크포인트, "가격정보_배포_모듈 launch")
	
	for i := 0; i < 구독_모듈_수량; i++ {
	    공용.F파이썬_프로세스_실행("price_data_test.py", "subscriber", 공용.P가격정보_배포_주소, 공용.P테스트_결과_회신_주소)
	    //공용.F체크포인트(&체크포인트, "파이썬 가격정보 구독 모듈 launch")
	}
	
	공용.F파이썬_프로세스_실행("price_data_test.py", "provider", 공용.P가격정보_입수_주소, strconv.Itoa(가격정보_배포횟수))
	//공용.F체크포인트(&체크포인트, "파이썬 가격정보 '제공' 모듈 launch")
	
	for i := 0; i < 구독_모듈_수량; i++ {
	    //공용.F체크포인트(&체크포인트, "테스트 결과 수신 RecvMessage() 시작", i)
	    메시지, 에러 := 테스트_결과_수신_소켓.RecvMessage(0)
	    //공용.F체크포인트(&체크포인트, "테스트 결과 수신 RecvMessage() 완료", i)

		if 에러 != nil {
			공용.F문자열_출력("테스트 결과 수신 중 에러 발생.\n %v\n %v\n", 에러.Error(), 공용.F변수_내역_문자열(메시지))
			
			테스트_결과_수신_소켓.SendMessage([]string{공용.P메시지_구분_에러, 에러.Error()})
			테스트.Fail()
		} else {
		    //공용.F체크포인트(&체크포인트, "테스트 결과 수신 후 회신 SendMessage() 시작", i)
		    테스트_결과_수신_소켓.SendMessage([]string{공용.P메시지_구분_OK, ""})
		    //공용.F체크포인트(&체크포인트, "테스트 결과 수신 후 회신 SendMessage() 완료", i)
		}
		
		구분 := 메시지[0]
		구독횟수 := 메시지[1]
		
		//공용.F체크포인트(&체크포인트, "결과 수신 반복문 테스트 시작", i)
		공용.F테스트_같음(테스트, 구분, 공용.P메시지_구분_일반)
		공용.F테스트_같음(테스트, 구독횟수, strconv.Itoa(가격정보_배포횟수))
		//공용.F체크포인트(&체크포인트, "결과 수신 반복문. 테스트 완료", i)
	}
	
	//공용.F체크포인트(&체크포인트, "테스트 종료")
}

func TestF가격정보_배포_모듈_Go(테스트 *testing.T) {
	공용.F멀티_스레드_모드()
	defer 공용.F단일_스레드_모드()

	가격정보_배포횟수 := 1000
	구독_모듈_수량 := 10
	결과값_채널_모음 := make([](chan int), 0)

	go F가격정보_배포_모듈()

	for i := 0; i < 구독_모듈_수량; i++ {
		공용.F초기화_대기열_추가(1)
		결과값_채널 := make(chan int)
		결과값_채널_모음 = append(결과값_채널_모음, 결과값_채널)

		go f테스트용_가격정보_구독_모듈(결과값_채널)
	}
	
	go f테스트용_가격정보_입수_모듈(가격정보_배포횟수)

	for i := 0; i < 구독_모듈_수량; i++ {
		결과값_채널 := 결과값_채널_모음[i]
		가격정보_구독횟수 := <-결과값_채널
		공용.F테스트_같음(테스트, 가격정보_구독횟수, 가격정보_배포횟수)
		//공용.F문자열_출력("완료 횟수 : %v", i+1)
	}
}

func f테스트용_가격정보_입수_모듈(가격정보_배포횟수 int) {
	// 가격정보_송신_소켓
	가격정보_송신_소켓, 에러 := zmq.NewSocket(zmq.REQ)
	defer 가격정보_송신_소켓.Close()

	if 에러 != nil {
		공용.F문자열_출력("가격정보_송신_소켓 초기화 중 에러 발생. %s", 에러.Error())
		panic(에러)
	}

	가격정보_송신_소켓.Connect(공용.P가격정보_입수_주소)

	//공용.F문자열_출력("f테스트용_가격정보_입수_모듈() 초기화 완료.")
	
	// 모든 모듈의 소켓이 안정화가 될 때까지 잠시 대기
	// 이러한 시간적 여유를 두지 않으면 구독 모듈에서 메시지 누락이 발생함.
	time.Sleep(time.Second)
	//공용.F문자열_출력("f테스트용_가격정보_입수_모듈() 실행 시작.")

	var 메시지 []string
	var 구분 string
	var 에러_메시지 string

	for i := 0; i < 가격정보_배포횟수; i++ {
		//공용.F문자열_출력("f테스트용_가격정보_입수_모듈() : 입수모듈 %v", i + 1)

		가격 := i * 10

		// 가격정보 송신
		메시지 = []string{공용.P메시지_구분_일반, strconv.Itoa(가격)}

		_, 에러 = 가격정보_송신_소켓.SendMessage(메시지)

		if 에러 != nil {
			공용.F문자열_출력("가격정보 송신 중 에러 발생.\n %v\n %v\n", 에러.Error(), 공용.F변수_내역_문자열(메시지[0], 메시지[1]))
			가격정보_송신_소켓.SendMessage([]string{공용.P메시지_구분_에러, 에러.Error()})
			//panic(에러)
			continue
		}

		//공용.F문자열_출력("f테스트용_가격정보_입수_모듈() : SendMessage %v", i + 1)

		메시지, 에러 = 가격정보_송신_소켓.RecvMessage(0)

		//공용.F문자열_출력("f테스트용_가격정보_입수_모듈() : RecvMessage %v", i + 1)

		if 에러 != nil {
			공용.F문자열_출력("가격정보 송신 후 회신 수신 중 에러 발생.\n %v\n %v\n", 에러.Error(), 공용.F변수_내역_문자열(메시지[0], 메시지[1]))
			//panic(에러)
			continue
		}

		구분 = 메시지[0]
		에러_메시지 = 메시지[1]

		if 구분 == 공용.P메시지_구분_에러 {
			공용.F문자열_출력("가격정보 송신 후 에러 메시지 수신.\n %v\n", 에러_메시지)
			//panic(에러_메시지)
			continue
		}
	}

	메시지 = []string{공용.P메시지_구분_종료, ""}
	가격정보_송신_소켓.SendMessage(메시지)
	가격정보_송신_소켓.RecvMessage(0)

	//공용.F문자열_출력("f테스트용_가격정보_입수_모듈() 종료.")
}

func f테스트용_가격정보_구독_모듈(결과값_채널 chan int) {
	가격정보_구독_소켓, 에러 := zmq.NewSocket(zmq.SUB)
	defer 가격정보_구독_소켓.Close()

	if 에러 != nil {
		공용.F문자열_출력("가격정보_구독_소켓 초기화 중 에러 발생. %v", 에러.Error())
		panic(에러)
	}

	가격정보_구독_소켓.Connect(공용.P가격정보_배포_주소)
	가격정보_구독_소켓.SetSubscribe("")

	//공용.F문자열_출력("f테스트용_가격정보_구독_모듈() 초기화 완료.")

	var 메시지 []string
	var 구분 string
	//var 데이터 string

	반복횟수 := 0
	가격정보_구독횟수 := 0

	for {
		메시지, 에러 = 가격정보_구독_소켓.RecvMessage(0)

		if 에러 != nil {
			공용.F문자열_출력("가격정보_구독_소켓 메시지 수신 중 에러 발생. %v", 에러.Error())
			continue
		}

		구분 = 메시지[0]
		//데이터 = 메시지[1]

		if 구분 == 공용.P메시지_구분_일반 {
			가격정보_구독횟수++
		} else if 구분 == 공용.P메시지_구분_종료 {
			break
		}

		반복횟수++
	}

	// 반복횟수는 0부터 시작했으나, 종료 메시지가 1회 포함되어 있으므로, 그대로 사용하면 됨.
	//공용.F문자열_출력("f테스트용_가격정보_구독_모듈() %v회 반복완료 후 종료.", 반복횟수)

	결과값_채널 <- 가격정보_구독횟수
}