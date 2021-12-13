var ss = require('simple-statistics');

var fs = require("fs");
var input = fs.readFileSync(0).toString();

data = input.trim().split('\n').map(d => +d);

let clusters = ss.ckmeans(data.slice(1), data[0]);

clusters.forEach((cluster, i) => {
    for (let d of cluster) {
        console.log(i, d);
    }
});
