const cont = document.getElementById("container")
const Http = new XMLHttpRequest();
let url = "http://localhost:8080/"
const SIZE = 5

let notepads = []

function init() {
    for (let i = 0; i < 1; i++) {
        notepads.push(new NotePad())
    }
    notepads.forEach(pad => {
        cont.appendChild(pad.elem)
    });
}

class NotePad {
    constructor() {
        this.elem = elemMaker("div", {
            classList: "notepad"
        })
        this.elem.addEventListener("click", e => {
            if (e.path[0].classList == "notepad") this.lines[this.lines.length - 1].focus()

        })
        this.elem.addEventListener("keydown", e => {
            if (e.key == "Enter") {
                e.preventDefault()
                console.log(e.path)
                this.addLine()
            }
        })
        this.lines = [new Line()]
        this.lines.forEach(e => {
            this.elem.appendChild(e.elem)
        })
    }
    addLine(i) {
        let line = new Line()
        this.lines.push(line)
        this.elem.appendChild(line.elem)
        line.focus()
    }
}

class Line {
    constructor() {
        this.elem = elemMaker("div", {
            classList: "line"
        })
        this.editPart = elemMaker("div", {
            contentEditable: "true",
            classList: "edit"
        })
        this.elem.appendChild(this.editPart)
        this.suggestor = elemMaker("span", {
            classList: "sug"
        })
        this.elem.appendChild(this.suggestor)
        this.editPart.addEventListener("keyup", e => {
            if (e.key != "Tab") {
                suggest(this.editPart.innerText, this.suggestor)
            }
        })
        this.editPart.addEventListener("keydown", e => {
            if (e.key == "Tab") {
                e.preventDefault()
                let selection = window.getSelection();
                let range = selection.getRangeAt(0);
                range.deleteContents();
                let node = document.createTextNode(this.suggestor.innerText);
                range.insertNode(node);
                selection.modify("move", "right", "paragraphboundary")
            }
        })
    }
    focus() {
        this.editPart.focus()
    }
}

function suggest(words, elem) {
    Http.open("GET", url + words.slice(words.length - SIZE, words.length));
    Http.send();
    Http.onreadystatechange = (e) => {
        elem.innerText = Http.responseText
    }
}

function elemMaker(elem, config) {
    let element = document.createElement(elem);
    for (key in config) element[key] = config[key];
    return element;
}

init()