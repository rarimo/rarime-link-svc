package data

type ProofOperator string

const (
	Noop ProofOperator = "$noop"
	Eq   ProofOperator = "$eq"
	Lt   ProofOperator = "$lt"
	Gt   ProofOperator = "$gt"
	In   ProofOperator = "$in"
	Nin  ProofOperator = "$nin"
	Ne   ProofOperator = "$ne"
)

func (o ProofOperator) String() string {
	return string(o)
}

var proofOperatorStrConv = map[string]ProofOperator{
	"$noop": Noop,
	"$eq":   Eq,
	"$lt":   Lt,
	"$gt":   Gt,
	"$in":   In,
	"$nin":  Nin,
	"$ne":   Ne,
}

func ProofOperatorFromString(data string) (ProofOperator, bool) {
	res, ok := proofOperatorStrConv[data]
	return res, ok
}

func MustProofOperatorFromString(data string) ProofOperator {
	return proofOperatorStrConv[data]
}
