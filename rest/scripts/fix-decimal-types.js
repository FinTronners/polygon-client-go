const fs = require('fs');
const path = require('path');
const { execFileSync } = require('child_process');

const filePath = path.join(__dirname, '../gen/client.gen.go');

let content = fs.readFileSync(filePath, 'utf8');

// Convert price/size fields from float64/float32 to decimal.Decimal
console.log('🔧 Converting price/size fields to decimal.Decimal...');

const decimalFields = ['AskPrice', 'BidPrice', 'AskSize', 'BidSize'];
let decimalCount = 0;

decimalFields.forEach(field => {
  const patterns = [
    [new RegExp(`(\\b${field}) (\\*?)float64`, 'g'), `$1 $2decimal.Decimal`],
    [new RegExp(`(\\b${field}) (\\*?)float32`, 'g'), `$1 $2decimal.Decimal`],
  ];
  patterns.forEach(([regex, repl]) => {
    const before = (content.match(regex) || []).length;
    content = content.replace(regex, repl);
    decimalCount += before;
  });
});

// Convert Price/Size fields in response structs only (exclude query params with form: tags)
const responsePriceFields = [
  [/(\bPrice) (\*?)float64 `json:/g, '$1 $2decimal.Decimal `json:'],
  [/(\bSize) (\*?)float64 `json:/g, '$1 $2decimal.Decimal `json:'],
  [/(\bPrice) (\*?)float32 `json:/g, '$1 $2decimal.Decimal `json:'],
  [/(\bSize) (\*?)float32 `json:/g, '$1 $2decimal.Decimal `json:'],
];
responsePriceFields.forEach(([regex, repl]) => {
  const before = (content.match(regex) || []).length;
  content = content.replace(regex, repl);
  decimalCount += before;
});

console.log(`✅ Decimal conversion: ${decimalCount} fields`);

// Inject decimal import if not already present
if (!content.includes('"github.com/shopspring/decimal"')) {
  content = content.replace(
    /("github.com\/oapi-codegen\/runtime\/types")/,
    '$1\n\t"github.com/shopspring/decimal"'
  );
  console.log('✅ Added decimal import');
}

fs.writeFileSync(filePath, content);

// Clean up formatting
try {
  execFileSync('gofmt', ['-w', filePath]);
  console.log('✅ Ran gofmt');
} catch (e) {
  console.log('⚠️  gofmt skipped (optional)');
}

console.log(`\n🎉 Decimal fix complete! ${decimalCount} fields converted.`);
