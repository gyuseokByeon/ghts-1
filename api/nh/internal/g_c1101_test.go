package internal

import (
	공용 "github.com/ghts/ghts/common"

	"math"
	"strings"
	"testing"
	"time"
	"unicode/utf8"
	"unsafe"
)

func TestCh조회_주식_현재가(테스트 *testing.T) {
	f접속_확인()

	종목 := 공용.F임의_종목_주식()
	질의 := 공용.New질의_가변형(P30초, 공용.P메시지_GET, TR주식_현재가_조회, 종목.G코드())
	질의.S질의(Ch조회)

	기본_자료 := new(S주식_현재가_조회_기본_자료)
	기본_자료 = nil

	변동_자료_모음 := make([]S주식_현재가_조회_변동_거래량_자료, 0)

	동시호가_자료 := new(S주식_현재가_조회_동시호가)
	동시호가_자료 = nil

	ok := true
	완료_메시지_수신 := false

	for !완료_메시지_수신 || 기본_자료 == nil {
		회신 := 질의.G회신()
		공용.F테스트_에러없음(테스트, 회신.G에러())

		switch 회신.G구분() {
		case P회신_조회:
			공용.F테스트_같음(테스트, 회신.G길이(), 1)

			수신_데이터, ok := 회신.G내용(0).(S수신_데이터)
			공용.F테스트_참임(테스트, ok)

			switch 수신_데이터.G블록_이름() {
			case "c1101OutBlock":
				공용.F테스트_같음(테스트, 수신_데이터.G길이(),
					int(unsafe.Sizeof(Tc1101OutBlock{})))
				공용.F테스트_다름(테스트, 수신_데이터.G데이터(), nil)

				기본_자료, ok = 수신_데이터.G데이터().(*S주식_현재가_조회_기본_자료)
				공용.F테스트_참임(테스트, ok)
			case "c1101OutBlock2":
				공용.F테스트_참임(테스트, 수신_데이터.G길이()%
					int(unsafe.Sizeof(Tc1101OutBlock2{})) == 0)
				공용.F테스트_다름(테스트, 수신_데이터.G데이터(), nil)

				변동_자료_모음, ok = 수신_데이터.G데이터().([]S주식_현재가_조회_변동_거래량_자료)
				공용.F테스트_참임(테스트, ok)
			case "c1101OutBlock3":
				공용.F테스트_같음(테스트, 수신_데이터.G길이(),
					int(unsafe.Sizeof(Tc1101OutBlock3{})))
				공용.F테스트_다름(테스트, 수신_데이터.G데이터(), nil)

				동시호가_자료, ok = 수신_데이터.G데이터().(*S주식_현재가_조회_동시호가)
				공용.F테스트_참임(테스트, ok)
			default:
				공용.F문자열_출력("예상치 못한 블록 이름 %v", 수신_데이터.G블록_이름())
				테스트.FailNow()
			}
		case P회신_메시지:
			공용.F테스트_같음(테스트, 회신.G길이(), 2)

			_, ok = 회신.G내용(0).(string) // 코드
			공용.F테스트_참임(테스트, ok)

			메시지, ok := 회신.G내용(1).(string)
			공용.F테스트_참임(테스트, ok)

			공용.F테스트_참임(테스트, strings.Contains(메시지, "조회완료"))
		case P회신_완료:
			공용.F테스트_같음(테스트, 회신.G길이(), 1)

			수신_데이터, ok := 회신.G내용(0).(S수신_데이터)
			공용.F테스트_참임(테스트, ok)
			공용.F테스트_같음(테스트, 수신_데이터.G블록_이름(), "c1101")
			공용.F테스트_같음(테스트, 수신_데이터.G길이(), 0)
			공용.F테스트_같음(테스트, 수신_데이터.G데이터(), nil)

			완료_메시지_수신 = true
		case P회신_에러:
			공용.F에러("P회신_에러 수신")
			테스트.FailNow()
		default:
			공용.F문자열_출력("\n*** %v 예상치 못한 회신 구분 : %v ***", 회신.G구분())
			공용.F변수값_확인(회신.G구분())
			공용.F변수값_확인(회신)
			테스트.FailNow()
		}
	}
	
	공용.F문자열_출력("*** 종목코드 %v ***", 종목.G코드())
	공용.F문자열_출력("*** 시각 %v ***", 기본_자료.M시각)
	
	// 기본 자료 테스트
	f주식_현재가_조회_기본_자료_테스트(테스트, 기본_자료, 종목)

	// 변동 자료 테스트
	f주식_현재가_조회_변동_거래량_자료_테스트(테스트, 기본_자료, 변동_자료_모음)

	// 동시호가 자료 테스트
	공용.F테스트_참임(테스트, 동시호가_자료 != nil, "동시호가 자료를 수신하지 못함.")
	f주식_현재가_조회_동시호가_자료_테스트(테스트, 기본_자료, 동시호가_자료)
}

func f주식_현재가_조회_기본_자료_테스트(테스트 *testing.T,
	s *S주식_현재가_조회_기본_자료, 종목 공용.I종목) {

	지금 := time.Now()
	삼분전 := 지금.Add(-3 * time.Minute)
	삼분후 := 지금.Add(3 * time.Minute)
	금일_0시 := time.Date(지금.Year(), 지금.Month(), 지금.Day(), 0, 0, 0, 0, 지금.Location())
	금일_9시 := 금일_0시.Add(9*time.Hour)
	개장_시각 := 금일_9시
	최근_개장일, 에러 := 공용.F한국증시_최근_개장일()
	공용.F에러_패닉(에러)
	
	삼십일전 := 지금.Add(-30 * 24 * time.Hour)
	연초 := time.Date(지금.Year(), time.January, 1, 0, 0, 0, 0, 지금.Location())	
	일년전 := 지금.Add(-1 * 366 * 24 * time.Hour)
	이백년전 := 지금.Add(-200 * 365 * 24 * time.Hour)

	공용.F테스트_참임(테스트, s != nil, "기본 자료를 수신하지 못함.")
	공용.F테스트_같음(테스트, s.M종목_코드, 종목.G코드())
	공용.F테스트_참임(테스트, utf8.ValidString(s.M종목명))
	공용.F테스트_참임(테스트, strings.Contains(s.M종목명, 종목.G이름()))
	공용.F테스트_참임(테스트, s.M등락율 >= 0) // 절대값임.

	f테스트_등락부호(테스트, s.M등락부호, s.M현재가, s.M전일_종가, s.M상한가, s.M하한가)
	공용.F테스트_같음(테스트, s.M전일_종가+f등락부호2정수(s.M등락부호)*s.M등락폭, s.M현재가)

	if s.M현재가 != 0 && s.M등락폭 != 0 && s.M등락율 != 0 {
		등락율_근사값 := math.Abs(float64(s.M등락폭)) / float64(s.M현재가) * 100
		공용.F테스트_참임(테스트, 공용.F오차율(등락율_근사값, s.M등락율) < 10)
	}

	공용.F테스트_참임(테스트, s.M거래량 >= 0)
	공용.F테스트_참임(테스트, s.M전일대비_거래량_비율 >= 0)

	거래량_비율_근사값 := float64(s.M거래량) / float64(s.M전일_거래량) * 100
	공용.F테스트_참임(테스트, 공용.F오차율(s.M전일대비_거래량_비율, 거래량_비율_근사값) < 10,
		s.M전일대비_거래량_비율, 거래량_비율_근사값, s.M거래량, s.M전일_거래량)

	if s.M유동_주식수_1000주 != 0 {
		유동주_회전율_근사값 := float64(s.M거래량) /
			float64(s.M유동_주식수_1000주*1000) * 100
		유동주_회전율_근사값 = math.Trunc(유동주_회전율_근사값*100) / 100
		공용.F테스트_참임(테스트, 공용.F오차(s.M유동주_회전율, 유동주_회전율_근사값) < 1 ||
			공용.F오차율(s.M유동주_회전율, 유동주_회전율_근사값) < 10,
			s.M유동주_회전율, 유동주_회전율_근사값)
	}

	if s.M거래대금_100만 != 0 && s.M거래량 != 0 && s.M현재가 != 0 {
		거래대금_근사값 := s.M거래량 * s.M현재가 / 1000000
		공용.F테스트_참임(테스트, 공용.F오차율(s.M거래대금_100만, 거래대금_근사값) < 10)
	}
	
	if s.M거래량 > 0 {
		// 거래량이 0이면 저가, 고가 모두 0임.
		공용.F테스트_참임(테스트, s.M저가 > 0, s.M저가)
		공용.F테스트_참임(테스트, s.M고가 > 0, s.M고가)
		
		공용.F테스트_참임(테스트, s.M저가 >= s.M하한가, s.M하한가, s.M저가)
		공용.F테스트_참임(테스트, s.M현재가 <= s.M고가, s.M현재가, s.M고가)
		공용.F테스트_참임(테스트, s.M상한가 >= s.M고가)
		공용.F테스트_참임(테스트, s.M고가 >= s.M시가)
		공용.F테스트_참임(테스트, s.M고가 >= s.M저가)
		공용.F테스트_참임(테스트, s.M시가 >= s.M저가)
		공용.F테스트_참임(테스트, s.M저가 >= s.M하한가)
		공용.F테스트_참임(테스트, s.M현재가 >= s.M저가)
		공용.F테스트_참임(테스트, s.M현재가 <= s.M고가)
		공용.F테스트_참임(테스트, s.M가중_평균_가격 >= s.M저가)
		공용.F테스트_참임(테스트, s.M가중_평균_가격 <= s.M고가)
	}

	공용.F테스트_참임(테스트, s.M하한가 > 0)
	공용.F테스트_참임(테스트, s.M연중_최저가 > 0)
	공용.F테스트_참임(테스트, s.M52주_고가 >= s.M연중_최고가)
	공용.F테스트_참임(테스트, s.M52주_고가 >= s.M20일_고가)
	공용.F테스트_참임(테스트, s.M20일_고가 >= s.M5일_고가)
	공용.F테스트_참임(테스트, s.M5일_고가 >= s.M5일_저가)
	공용.F테스트_참임(테스트, s.M연중_최저가 >= s.M52주_저가)
	공용.F테스트_참임(테스트, s.M20일_저가 >= s.M52주_저가)
	공용.F테스트_참임(테스트, s.M5일_저가 >= s.M20일_저가)
	공용.F테스트_참임(테스트, s.M연중_최고가 >= s.M연중_최저가)
	f테스트_등락부호(테스트, s.M시가대비_등락부호, s.M현재가, s.M시가, s.M상한가, s.M하한가)
	공용.F테스트_같음(테스트, s.M시가+s.M시가대비_등락폭, s.M현재가) // 시가대비_등락폭 자체에 부호가 반영되어 있음.
	공용.F테스트_참임(테스트, s.M시각.After(최근_개장일.Add(-1*time.Second)))
	공용.F테스트_참임(테스트, s.M시각.Before(최근_개장일.Add(16*time.Hour)))
	공용.F테스트_참임(테스트, s.M시각.Before(삼분후))

	if 공용.F한국증시_장중() { // 장중
		공용.F테스트_참임(테스트, s.M시각.After(삼분전), s.M시각, 삼분전)
		공용.F테스트_참임(테스트, s.M시각.Before(삼분후), s.M시각, 삼분후)
	} else { // 장중이 아니면 마감 시각 기록.		
		공용.F테스트_참임(테스트, s.M시각.Hour() == 15, s.M시각)
	}

	매도_잔량_합계 := int64(0)
	for i, 매도_잔량 := range s.M매도_잔량_모음 {
		공용.F테스트_참임(테스트, 매도_잔량 >= 0, i, 매도_잔량)
		
		if 매도_잔량 == 0 {
			continue
		}
		
		매도_잔량_합계 += 매도_잔량
		매도_호가 := s.M매도_호가_모음[i]
		공용.F테스트_참임(테스트, 매도_호가 <= s.M상한가)
		공용.F테스트_참임(테스트, 매도_호가 >= s.M하한가)
		
		switch i {
		case 0:
			공용.F테스트_참임(테스트, 매도_호가 >= s.M현재가)
		default:
			공용.F테스트_참임(테스트, 매도_호가 > s.M매도_호가_모음[i-1])
		} 
	}
	
	매수_잔량_합계 := int64(0)
	for i, 매수_잔량 := range s.M매수_잔량_모음 {
		공용.F테스트_참임(테스트, 매수_잔량 >= 0, i, 매수_잔량)
		
		if 매수_잔량 == 0 {
			continue
		}
		
		매수_호가 := s.M매수_호가_모음[i]
		공용.F테스트_참임(테스트, 매수_호가 <= s.M상한가)
		공용.F테스트_참임(테스트, 매수_호가 >= s.M하한가)
		
		switch i {
		case 0:
			공용.F테스트_참임(테스트, 매수_호가 <= s.M현재가)
		default:
			공용.F테스트_참임(테스트, 매수_호가 < s.M매수_호가_모음[i-1],
				i, 매수_호가, s.M매수_호가_모음[i-1])
		}
	}

	공용.F테스트_참임(테스트, s.M매도_잔량_총합 >= 매도_잔량_합계)
	공용.F테스트_참임(테스트, s.M매수_잔량_총합 >= 매수_잔량_합계)
	공용.F테스트_참임(테스트, s.M시간외_매도_잔량 >= 0)
	공용.F테스트_참임(테스트, s.M시간외_매수_잔량 >= 0)
	공용.F테스트_참임(테스트, s.M피봇_2차_저항 >= s.M피봇_1차_저항)
	공용.F테스트_참임(테스트, s.M피봇_1차_저항 >= s.M피봇가)
	공용.F테스트_참임(테스트, s.M피봇가 >= s.M피봇_1차_지지)
	공용.F테스트_참임(테스트, s.M피봇_1차_지지 >= s.M피봇_2차_지지)
	공용.F테스트_참임(테스트, utf8.ValidString(s.M시장_구분))
	공용.F테스트_같음(테스트, s.M시장_구분, "코스피", "코스닥")
	공용.F테스트_참임(테스트, utf8.ValidString(s.M업종명))
	공용.F테스트_참임(테스트, utf8.ValidString(s.M자본금_규모))
	공용.F테스트_참임(테스트, strings.Contains(s.M자본금_규모, "형주"))
	공용.F테스트_참임(테스트, utf8.ValidString(s.M결산월))
	공용.F테스트_참임(테스트, strings.Contains(s.M결산월, "월 결산"))
	
	for _, 추가_정보 := range s.M추가_정보_모음 {
		공용.F테스트_참임(테스트, utf8.ValidString(추가_정보))
	}

	공용.F테스트_참임(테스트, utf8.ValidString(s.M서킷_브레이커_구분))
	공용.F테스트_같음(테스트, s.M서킷_브레이커_구분, "", "CB발동", "CB해제", "장종료")
	공용.F테스트_참임(테스트, s.M액면가 > 0)
	//공용.F테스트_참임(테스트, strings.Contains(s.M전일_종가_타이틀, "전일종가"))
	공용.F테스트_참임(테스트, 공용.F오차율(s.M상한가, float64(s.M전일_종가)*1.3) < 5)
	공용.F테스트_참임(테스트, 공용.F오차율(s.M하한가, float64(s.M전일_종가)*0.7) < 5)
	공용.F테스트_참임(테스트, s.M대용가 < s.M전일_종가)
	공용.F테스트_참임(테스트, s.M대용가 > int64(float64(s.M전일_종가)*0.5))
	공용.F테스트_참임(테스트, s.M공모가 >= 0)
	공용.F테스트_참임(테스트, s.M52주_저가_일자.After(일년전), s.M52주_저가_일자)
	공용.F테스트_참임(테스트, s.M52주_저가_일자.Before(지금), s.M52주_저가_일자)
	공용.F테스트_참임(테스트, s.M52주_고가_일자.After(일년전), s.M52주_고가_일자)
	공용.F테스트_참임(테스트, s.M52주_고가_일자.Before(지금), s.M52주_고가_일자)
	//공용.F테스트_참임(테스트, 공용.F오차(s.M상장_주식수_1000주 - (s.M상장_주식수/1000) <= 1.01 ||
	//	공용.F오차율(s.M상장_주식수_1000주 - (s.M상장_주식수/1000)) < 10)
	공용.F테스트_참임(테스트, s.M유동_주식수_1000주 >= 0)

	시가총액_근사값 := s.M현재가 * s.M상장_주식수 / 100000000
	공용.F테스트_참임(테스트, 공용.F오차율(s.M시가_총액_억, 시가총액_근사값) < 10)
	공용.F테스트_참임(테스트, s.M거래원_정보_수신_시각.Before(삼분후), 
		s.M거래원_정보_수신_시각, 삼분후)
	
	if 공용.F한국증시_장중() {
		공용.F테스트_참임(테스트, s.M거래원_정보_수신_시각.After(개장_시각))
		공용.F테스트_참임(테스트, s.M시각.Before(삼분후))
	} else {
		공용.F테스트_참임(테스트, s.M거래원_정보_수신_시각.Hour() == 15)
	}
	
	매도_거래량_합계 := int64(0)
	for i, 매도_거래량 := range s.M매도_거래량_모음 {
		공용.F테스트_참임(테스트, 매도_거래량 >= 0, i, 매도_거래량)
		
		if 매도_거래량 == 0 {
			continue
		}
		
		매도_거래량_합계 += 매도_거래량
		매도_거래원 := s.M매도_거래원_모음[i]
		공용.F테스트_참임(테스트, len(매도_거래원) > 0)
		공용.F테스트_참임(테스트, utf8.ValidString(매도_거래원), 매도_거래원)	
	}
	
	매수_거래량_합계 := int64(0)
	for i, 매수_거래량 := range s.M매수_거래량_모음 {
		공용.F테스트_참임(테스트, 매수_거래량 >= 0, i, 매수_거래량)
		
		if 매수_거래량 == 0 {
			continue
		}
		
		매수_거래량_합계 += 매수_거래량
		매수_거래원 := s.M매수_거래원_모음[i]
		공용.F테스트_참임(테스트, len(매수_거래원) > 0)
		공용.F테스트_참임(테스트, utf8.ValidString(매수_거래원), 매수_거래원)
	}
	
	공용.F테스트_참임(테스트, s.M외국인_매도_거래량 >= 0)
	공용.F테스트_참임(테스트, s.M외국인_매수_거래량 >= 0)
	공용.F테스트_참임(테스트, s.M외국인_시간.After(최근_개장일),
		s.M외국인_시간, 최근_개장일)
	공용.F테스트_참임(테스트, s.M외국인_시간.Before(최근_개장일.Add(23*time.Hour)),
		s.M외국인_시간, 최근_개장일.Add(23*time.Hour))
	
	if 공용.F한국증시_장중() {
		공용.F테스트_참임(테스트, s.M외국인_시간.Before(삼분후))
	}
	
	공용.F테스트_참임(테스트, s.M외국인_지분율 >= 0)
	공용.F테스트_참임(테스트, s.M외국인_지분율 <= 100)
	공용.F테스트_참임(테스트, s.M신용잔고_기준_결제일.After(삼십일전))
	공용.F테스트_참임(테스트, s.M신용잔고_기준_결제일.Before(최근_개장일.Add(16*time.Hour)))
	
	if 공용.F한국증시_장중() {
		공용.F테스트_참임(테스트, s.M신용잔고_기준_결제일.Before(금일_0시))
	}
		
	공용.F테스트_참임(테스트, s.M신용잔고율 >= 0)
	공용.F테스트_참임(테스트, s.M신용잔고율 <= 100)
	//공용.F테스트_참임(테스트, s.M유상_기준일.After(이백년전) || s.M유상_기준일.IsZero())
	//공용.F테스트_참임(테스트, s.M무상_기준일.After(이백년전) || s.M무상_기준일.IsZero())
	공용.F테스트_참임(테스트, s.M유상_배정_비율 >= 0)
	공용.F테스트_참임(테스트, s.M유상_배정_비율 <= 100)
	//공용.F테스트_참임(테스트, s.M외국인_순매수량 >= 0, s.M외국인_순매수량)	// 순매도 시 (-) 값을 가질 수 있음.
	공용.F테스트_참임(테스트, s.M무상_배정_비율 >= 0)
	공용.F테스트_참임(테스트, s.M무상_배정_비율 <= 100)
	//공용.F변수값_확인(s.M당일_자사주_신청_여부)

	공용.F테스트_참임(테스트, s.M상장일.After(이백년전))
	공용.F테스트_참임(테스트, s.M상장일.Before(최근_개장일.Add(1*time.Second)))
	공용.F테스트_참임(테스트, s.M대주주_지분율 >= 0)
	공용.F테스트_참임(테스트, s.M대주주_지분율 <= 100)
	공용.F테스트_참임(테스트, s.M대주주_지분율_정보_일자.After(일년전))
	공용.F테스트_참임(테스트, s.M대주주_지분율_정보_일자.Before(삼분후))
	//공용.F변수값_확인(s.M네잎클로버_종목_여부)	// NH투자증권 선정 추천 종목
	공용.F테스트_참임(테스트, s.M증거금_비율 >= 0)
	공용.F테스트_참임(테스트, s.M증거금_비율 <= 100)
	공용.F테스트_참임(테스트, s.M자본금 > 0)
	공용.F테스트_참임(테스트, s.M전체_거래원_매도_합계 >= 매도_거래량_합계)
	공용.F테스트_참임(테스트, s.M전체_거래원_매수_합계 >= 매수_거래량_합계)

	//공용.F변수값_확인(s.M종목명2)
	//공용.F테스트_참임(테스트, utf8.ValidString(s.M종목명2))
	//공용.F변수값_확인(s.M우회_상장_여부)	// 이 항목은 뭐하는 데 필요할까?

	//공용.F테스트_참임(테스트, s.M코스피_구분_2 == "코스피" || s.M코스피_구분_2 == "코스닥")
	//공용.F테스트_참임(테스트, utf8.ValidString(s.M코스피_구분_2))   // 앞에 나온 '코스피/코스닥 구분'과 중복 아닌가?
	
	공용.F테스트_참임(테스트, s.M공여율_기준일.After(삼십일전))
	공용.F테스트_참임(테스트, s.M공여율_기준일.Before(지금)) // 공여율은 '신용거래 관련 비율'이라고 함.
	공용.F테스트_참임(테스트, s.M공여율 >= 0 && s.M공여율 <= 100) // 공여율(%)
	공용.F테스트_참임(테스트, math.Abs(float64(s.PER)) < 100)
	공용.F테스트_참임(테스트, s.M종목별_신용한도 >= 0)
	공용.F테스트_참임(테스트, s.M종목별_신용한도 <= 100)
	공용.F테스트_참임(테스트, s.M가중_평균_가격 >= s.M저가)
	공용.F테스트_참임(테스트, s.M가중_평균_가격 <= s.M고가)
	공용.F테스트_참임(테스트, s.M추가_상장_주식수 >= 0)
	공용.F테스트_참임(테스트, utf8.ValidString(s.M종목_코멘트))
	공용.F테스트_참임(테스트, s.M전일_거래량 >= 0)
	공용.F테스트_참임(테스트, s.M전일_등락폭 >= 0) // 절대값
	공용.F테스트_참임(테스트, f올바른_등락부호(s.M전일_등락부호))
	공용.F테스트_참임(테스트, s.M연중_최고가_일자.After(연초))
	공용.F테스트_참임(테스트, s.M연중_최고가_일자.Before(지금))
	공용.F테스트_참임(테스트, s.M연중_최저가_일자.After(연초))
	공용.F테스트_참임(테스트, s.M연중_최저가_일자.Before(지금))
	공용.F테스트_참임(테스트, s.M외국인_보유_주식수 <= s.M상장_주식수+s.M추가_상장_주식수)
	공용.F테스트_참임(테스트, s.M외국인_지분_한도 >= 0)
	공용.F테스트_참임(테스트, s.M외국인_지분_한도 <= 100)
	공용.F테스트_참임(테스트, s.M매매_수량_단위 == 1, s.M매매_수량_단위)
	공용.F테스트_같음(테스트, int(s.M대량_매매_방향), 0, 1, 2)
	
	if s.M대량_매매_방향 == 0 {
		공용.F테스트_거짓임(테스트, s.M대량_매매_존재)
	} else {
		공용.F테스트_참임(테스트, s.M대량_매매_존재)
	}
}

func f주식_현재가_조회_변동_거래량_자료_테스트(테스트 *testing.T,
	기본_자료 *S주식_현재가_조회_기본_자료,
	변동_자료_모음 []S주식_현재가_조회_변동_거래량_자료) {
	공용.F테스트_참임(테스트, len(변동_자료_모음) > 0, "변동 자료를 수신하지 못함.")

	거래량_잔량 := 기본_자료.M거래량
	지금 := time.Now()
	삼분후 := 지금.Add(3 * time.Minute)
	최근_개장일, 에러 := 공용.F한국증시_최근_개장일()
	공용.F에러_패닉(에러)
	
	for i, s := range 변동_자료_모음 {
		공용.F테스트_참임(테스트, s.M시각.After(최근_개장일.Add(9*time.Hour)))
		공용.F테스트_참임(테스트, s.M시각.Before(삼분후), s.M시각)
		공용.F테스트_참임(테스트, s.M매도_호가 >= 0)
		공용.F테스트_참임(테스트, s.M매수_호가 >= 0)
		공용.F테스트_참임(테스트, s.M매도_호가 >= 기본_자료.M하한가 ||
			s.M매도_호가 == 0)
		공용.F테스트_참임(테스트, s.M매도_호가 <= 기본_자료.M상한가)
		공용.F테스트_참임(테스트, s.M매수_호가 >= 기본_자료.M하한가 ||
			s.M매수_호가 == 0)
		공용.F테스트_참임(테스트, s.M매수_호가 <= 기본_자료.M상한가)
		공용.F테스트_참임(테스트, s.M현재가 <= 기본_자료.M상한가)
		공용.F테스트_참임(테스트, s.M현재가 >= 기본_자료.M하한가)
		
		if 공용.F한국증시_장중() {
			공용.F테스트_참임(테스트, s.M시각.Before(삼분후), s.M시각)
			공용.F테스트_참임(테스트, s.M매도_호가 >= s.M현재가 ||
				s.M매도_호가 == 0)
			공용.F테스트_참임(테스트, s.M매수_호가 <= s.M현재가 ||
				s.M매수_호가 == 0)
		} else {
			공용.F테스트_참임(테스트, s.M시각.Before(최근_개장일.Add(16*time.Hour)))
			
			// 장 마감 후 매도호가, 매수호가는 흔히 생각하는 조건을 만족시키지 않음.
			//공용.F테스트_참임(테스트, s.M매도_호가 >= s.M현재가)
			//공용.F테스트_참임(테스트, s.M매수_호가 <= s.M현재가)						
		}

		공용.F테스트_참임(테스트, f올바른_등락부호(s.M등락부호))
		공용.F테스트_같음(테스트, f등락부호2정수(s.M등락부호)*s.M등락폭,
			s.M현재가-기본_자료.M전일_종가)
		
		// 걸러낸 자료로 인한 오차 수정.
		if i == 0 && s.M거래량 != 거래량_잔량 {
			거래량_잔량 = s.M거래량
		}

		공용.F테스트_같음(테스트, s.M거래량, 거래량_잔량)
		거래량_잔량 -= s.M변동_거래량
	}
}

func f주식_현재가_조회_동시호가_자료_테스트(테스트 *testing.T,
	기본_자료 *S주식_현재가_조회_기본_자료,
	s *S주식_현재가_조회_동시호가) {
	공용.F테스트_다름(테스트, s, nil)
	공용.F테스트_같음(테스트, int(s.M동시호가_구분), 0,1,2,3,4,5,6)

	if s.M동시호가_구분 == 0 { // 동시호가 아님.
		return
	}

	공용.F변수값_확인(기본_자료.M시각, 기본_자료.M종목_코드, s.M동시호가_구분)
	공용.F테스트_참임(테스트, f올바른_등락부호(s.M예상_체결_부호), s.M예상_체결_부호)
	공용.F테스트_참임(테스트, s.M예상_체결가 <= 기본_자료.M상한가)
	공용.F테스트_참임(테스트, s.M예상_체결가 >= 기본_자료.M하한가)
	공용.F테스트_참임(테스트, 공용.F오차율(s.M예상_체결가, 기본_자료.M현재가) < 10)
	공용.F테스트_같음(테스트, f등락부호2정수(s.M예상_체결_부호)*s.M예상_등락폭,
		s.M예상_체결가-기본_자료.M전일_종가)

	if s.M예상_등락폭 != 0 && s.M예상_등락율 != 0 {
		예상_등락율_근사값 := math.Abs(float64(s.M예상_등락폭)) /
			float64(s.M예상_체결가) * 100
		공용.F테스트_참임(테스트, 공용.F오차율(s.M예상_등락율, 예상_등락율_근사값) < 10)
	}

	공용.F테스트_참임(테스트, s.M예상_체결_수량 >= 0)
	공용.F테스트_참임(테스트, s.M예상_체결_수량 <= 기본_자료.M매도_잔량_총합 ||
		s.M예상_체결_수량 <= 기본_자료.M매수_잔량_총합)
}