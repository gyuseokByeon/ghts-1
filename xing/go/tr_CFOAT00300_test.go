/* Copyright (C) 2015-2019 김운하(UnHa Kim)  < unha.kim.ghts at gmail dot com >

이 파일은 GHTS의 일부입니다.

이 프로그램은 자유 소프트웨어입니다.
소프트웨어의 피양도자는 자유 소프트웨어 재단이 공표한 GNU LGPL 2.1판
규정에 따라 프로그램을 개작하거나 재배포할 수 있습니다.

이 프로그램은 유용하게 사용될 수 있으리라는 희망에서 배포되고 있지만,
특정한 목적에 적합하다거나, 이익을 안겨줄 수 있다는 묵시적인 보증을 포함한
어떠한 형태의 보증도 제공하지 않습니다.
보다 자세한 사항에 대해서는 GNU LGPL 2.1판을 참고하시기 바랍니다.
GNU LGPL 2.1판은 이 프로그램과 함께 제공됩니다.
만약, 이 문서가 누락되어 있다면 자유 소프트웨어 재단으로 문의하시기 바랍니다.
(자유 소프트웨어 재단 : Free Software Foundation, Inc.,
59 Temple Place - Suite 330, Boston, MA 02111-1307, USA)

Copyright (C) 2015-2019년 UnHa Kim (< unha.kim.ghts at gmail dot com >)

This file is part of GHTS.

GHTS is free software: you can redistribute it and/or modify
it under the terms of the GNU Lesser General Public License as published by
the Free Software Foundation, version 2.1 of the License.

GHTS is distributed in the hope that it will be useful,
but WITHOUT ANY WARRANTY; without even the implied warranty of
MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
GNU Lesser General Public License for more details.

You should have received a copy of the GNU Lesser General Public License
along with GHTS.  If not, see <http://www.gnu.org/licenses/>. */

package xing

import (
	"github.com/ghts/ghts/lib"
	xt "github.com/ghts/ghts/xing/base"
	"testing"
)

func TestCFOAT00300_선물옵션_취소주문_질의값(t *testing.T) {
	_, ok := interface{}(new(xt.CFOAT00300_선물옵션_취소주문_질의값)).(lib.I질의값)
	lib.F테스트_참임(t, ok)
}

func TestCFOAT00300_선물옵션_취소주문(t *testing.T) {
	t.Parallel()

	t.SkipNow()

	if !F한국증시_정규시장_거래시간임() {
		t.SkipNow()
	}

	계좌번호 := lib.F확인(F계좌_번호(1)).(string)	// 선물옵션 계좌를 선택해야 함.
	계좌_상세명, 에러 := F계좌_상세명(계좌번호)
	lib.F테스트_에러없음(t, 에러)
	lib.F테스트_같음(t, 계좌_상세명, "선물옵션")

	종목코드 := "??"
	취소수량 := int64(lib.F임의_범위_이내_정수값(0, 3))
	원주문번호 := int64(0)

	응답값, 에러 := TrCFOAT00300_선물옵션_취소주문(종목코드, 계좌번호, 원주문번호, 취소수량)
	lib.F테스트_에러없음(t, 에러)

	값1 := 응답값.M응답1
	lib.F테스트_같음(t, 값1.M레코드갯수, 1)
	lib.F테스트_같음(t, 값1.M주문시장, xt.P주문시장_비상장, xt.P주문시장_코스피, xt.P주문시장_채권,
		xt.P주문시장_장외시장, xt.P주문시장_코스닥, xt.P주문시장_코넥스, xt.P주문시장_프리보드,
		xt.P주문시장_동경거래소, xt.P주문시장_JASDAQ)
	lib.F테스트_같음(t, 값1.M계좌번호, 계좌번호)
	lib.F테스트_같음(t, 값1.M종목코드, 종목코드)
	lib.F테스트_같음(t, 값1.M주문유형, xt.P주문유형_해당없음, xt.P주문_현금매도, xt.P주문_현금매수,
		xt.P주문_신용매도, xt.P주문_신용매수, xt.P주문_저축매도, xt.P주문_저축매수, xt.P주문_상품매도_대차,
		xt.P주문_상품매도, xt.P주문_상품매수, xt.P주문_선물대용매도_일반, xt.P주문_선물대용매도_반대,
		xt.P주문_현금매도_프, xt.P주문_현금매수_프, xt.P주문_현금매수_유가, xt.P주문_현금매수_정리,
		xt.P주문_상품매도_대차_프, xt.P주문_상품매도_프, xt.P주문_상품매수_프, xt.P주문_장외매매)
	lib.F테스트_같음(t, 값1.M원주문번호, 원주문번호)
	lib.F테스트_참임(t, 값1.M취소수량 > 0, 값1.M취소수량)
	lib.F테스트_같음(t, 값1.M통신매체, xt.P통신매체_아이폰, xt.P통신매체_안드로이드,
		xt.P통신매체_API, xt.P통신매체_HTS, xt.P통신매체_모의서버_HTS)
	lib.F테스트_참임(t, lib.F2일자(값1.M협의매매완료시각).Equal(lib.F금일()), 값1.M협의매매완료시각)
	//M그룹ID     T증권그룹
	lib.F테스트_참임(t, 값1.M주문번호 > 0, 값1.M주문번호)
	lib.F체크포인트(값1.M포트폴리오번호)
	lib.F체크포인트(값1.M바스켓번호)
	lib.F체크포인트(값1.M트렌치번호)
	lib.F체크포인트(값1.M항목번호)
	lib.F체크포인트(값1.M관리사원번호)
	lib.F체크포인트(값1.M펀드ID)
	lib.F체크포인트(값1.M펀드원주문번호)
	lib.F체크포인트(값1.M펀드주문번호)

	값2 := 응답값.M응답2
	lib.F테스트_같음(t, 값2.M레코드갯수, 1)
	lib.F테스트_참임(t, 값2.M주문번호 > 0, 값2.M주문번호)
	lib.F체크포인트(값2.M지점명)
	lib.F체크포인트(값2.M계좌명)
	lib.F테스트_다름(t, 값2.M종목명, "")
	lib.F테스트_참임(t, 값2.M주문가능금액 > 0, 값2.M주문가능금액)
	lib.F테스트_참임(t, 값2.M현금주문가능금액 > 0, 값2.M현금주문가능금액)
	lib.F테스트_참임(t, 값2.M주문증거금액 > 0, 값2.M주문증거금액)
	lib.F테스트_참임(t, 값2.M현금주문증거금액 > 0, 값2.M현금주문증거금액)
	lib.F테스트_참임(t, 값2.M주문가능수량 > 0, 값2.M주문가능수량)
}
