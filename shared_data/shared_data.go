package shared_data

import (
	공용 "github.com/ghts/ghts/shared"
	
	"strconv"
	"sync"
)

var Ch주소 = make(chan 공용.I질의, 100)
var Ch종목 = make(chan 공용.I질의, 100)

var 공용_데이터_Go루틴_실행_중 = false
var 공용_데이터_Go루틴_잠금 = &sync.RWMutex{}

func F공용_데이터_Go루틴_실행_중() bool {
	공용_데이터_Go루틴_잠금.RLock()
	defer 공용_데이터_Go루틴_잠금.RUnlock()

	return 공용_데이터_Go루틴_실행_중
}

func F공용_데이터_Go루틴(Go루틴_생성_결과 chan bool) {
	if F공용_데이터_Go루틴_실행_중() {
		Go루틴_생성_결과 <- false
		return
	}
	
	공용_데이터_Go루틴_잠금.Lock()
	
	if 공용_데이터_Go루틴_실행_중 {
		Go루틴_생성_결과 <- false
		공용_데이터_Go루틴_잠금.Unlock()
		return
	}
	
	공용_데이터_Go루틴_실행_중 = true
	공용_데이터_Go루틴_잠금.Unlock()
	
	주소_맵 := f_주소_맵_초기화()
	종목_맵 := f종목_맵_초기화()
	
	// 초기화 완료
	Go루틴_생성_결과 <- true
	
	for {
		select {
		case 질의 := <-Ch주소:
			var 회신 공용.I회신		
			
			switch {
			case 질의.G구분() != 공용.P메시지_일반:
				회신 = 공용.New회신(
							공용.F에러_생성("잘못된 메시지 구분 '%v'", 질의.G구분()),
							공용.P메시지_에러)	
			case 질의.G길이() != 1:
				회신 = 공용.New회신(
							공용.F에러_생성("잘못된 질의 내용 길이 %v", 질의.G길이()),
							공용.P메시지_에러)
			default:
				주소, 존재함 := 주소_맵[질의.G내용(0)]
				
				if !존재함 {
					회신 = 공용.New회신(
								공용.F에러_생성("잘못된 질의값 '%v'", 질의.G내용(0)),
								공용.P메시지_에러)
				} else {
					회신 = 공용.New회신(nil, 공용.P메시지_OK, 주소)
				}
			}
			
			질의.G회신_채널() <- 회신
		case 질의 := <-Ch종목:
			var 회신 공용.I회신		
			
			switch {
			case 질의.G구분() != 공용.P메시지_일반:
				회신 = 공용.New회신(
							공용.F에러_생성("잘못된 메시지 구분 '%v'", 질의.G구분()),
							공용.P메시지_에러)
			case 질의.G길이() != 1:
				회신 = 공용.New회신(
							공용.F에러_생성("잘못된 질의 내용 길이 %v", 질의.G길이()),
							공용.P메시지_에러)
			default:
				종목, 존재함 := 종목_맵[질의.G내용(0)]
				
				if !존재함 {
					회신 = 공용.New회신(
								공용.F에러_생성("잘못된 질의값 '%v', %v", 질의.G내용(0), 주소_맵),
								공용.P메시지_에러)
				} else {
					회신 = 공용.New회신(nil, 공용.P메시지_OK, 종목.G코드(), 종목.G이름())
				}
			}
			
			질의.G회신_채널() <- 회신
		}
	}
}

func f_주소_맵_초기화() map[string]string {
	맵 := make(map[string]string)
	
	맵[P주소정보] = 공용.P주소_주소정보
	맵[P테스트_결과] = 공용.P주소_테스트_결과
	
	주소_모음 := make([]string, 0)
	주소_모음 = append(주소_모음, P종목정보)
	주소_모음 = append(주소_모음, P가격정보)
	주소_모음 = append(주소_모음, P가격정보_입수)
	주소_모음 = append(주소_모음, P가격정보_배포) 

	for i:=0 ; i < len(주소_모음) ; i++ {
		맵[주소_모음[i]] = "tcp://127.0.0.1:" + strconv.Itoa(3010 + i)	// 3010번 포트부터 차례대로 배정.	
	}
	
	return 맵
}

func f종목_맵_초기화() map[string]공용.I종목 {
	맵 := make(map[string]공용.I종목)
	
	공용.F메모("f종목_맵_초기화() 개발할 것.")
	
	// 임시로 샘플 데이터만 사용해서 테스트 할 수 있도록 함.
	종목_모음 := 공용.F샘플_종목_모음()

	for i:=0 ; i < len(종목_모음) ; i++ {
		맵[종목_모음[i].G코드()] = 종목_모음[i]
	}
	
	return 맵
}