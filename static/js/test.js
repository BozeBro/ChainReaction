let t = (h = 1) => {
    if (h == 10) return 10
    console.log(h)
    if (h + 1 === 10) {console.log(true)}
    return new Promise(() => t(h+1))
}
let n = () => new Promise(() => t());
let m = async () => {
    console.log("YO");
    await n();
    console.log("NO")
}
for (let i = 0; i < 10; i++) {
    m()
    console.log("HfI")
}