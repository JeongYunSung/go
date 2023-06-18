# protobuf

프로토콜버퍼는 구조화된 데이터를 직렬화하기 위한 언어 중립적인 메시지 프로토콜이다.

## benefits

1. compact한 data storage
2. 빠른 구문 분석
3. 다양한 언어 지원
4. 자동 생성을 통한 최적화된 기능

## Warring

1. protobuf는 메시지 전체를 메모리에 올려두는데 이에 몇 메가바이트를 초과하는 경우 메모리 사용량이 급증할 수 있음
2. 동일한 데이터가 직렬화되면 다양한 이진 직렬화를 가질 수 있으므로 동일하지 않음. 따라서 항상 parsing을 완수한 후 비교를 해야 됨
3. 메시지는 기본적으로 압축되지 않음. gzip과 같은 압축 알고리즘을 이용
4. 부동소수점, 다차원 배열 등을 사용하는 과학 및 엔지니어링 용도로는 크기 및 속도 효율이 최대치를 기대할 순 없음
5. `proto`라는 스키마 파일이 있음

## Who

1. gRPC
2. GCM

# 메시지 유형 정의

```protobuf
syntax = "proto3";

message SearchRequest {
  string query = 1;
  int32 page_number = 2;
  int32 result_per_page = 3;
}
```
1. syntax를 사용해 protobuf 버전지정
2. 스칼라 및 복합 필드를 지정
3. 필드 번호 할당
   1. 주어진 번호는 고유해야 함
   2. 19000~19999는 예약된 번호로 사용할 수 없음
   3. 1~15의 경우 1byte만 사용 함 ( 16~2047은 2바이트 )
   4. 메시지 유형이 사용중이면 메시지 번호는 변경할 수 없음

## 필드 레이블

1. optional : 값이 있거나 기본값을 반환하거나
2. repeated : 0번이상 반복될 수 있음
3. map : key value 필드 유형
4. 암시적 필드 존재

## 필드 삭제

`reserved 2, 5, 7`과 같이 해당 필드는 삭제되었다고 명시해주어야 함

## 동작 방식

1. proto파일에 schema 작성
2. protoc 컴파일러를 사용해 code generator ( .java, .py, .cc, .go 등의 확장자 파일로 만들어 줌 )
3. 위에서 생성된 파일을 실제 프로젝트 코드에서 사용
4. 라이브러리나 gRPC등을 이용해 serializable 및 deserializable