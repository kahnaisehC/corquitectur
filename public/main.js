
let fileInput = document.getElementById("fileInput")
let selectDiv = document.getElementById("selectDiv")
console.log(fileInput.attributes)
fileInput.addEventListener("input", function(){
    let file = fileInput.files["0"]
    fileInput.files["0"]
    .stream()
    .getReader()
    .read()
    .then((r) => {
        let cutCharCode = Number('\n'.charCodeAt(0))
        let headersString = ""
        for(let i = 0; r.value[i] != cutCharCode; i++){
            headersString += String.fromCharCode(r.value[i])
        }
        let headers = headersString.split(",")
        for(let i = 0; i < headers.length; i++){
            let header = headers[i]
            textInput = document.createElement("input")
            select = document.createElement("select")
            let opt = document.createElement("option")
            console.log(select.attributes)
            opt.value = "string"
            opt.innerText = "string"
            select.name = "select"+header
            select.append(opt)
            textInput.value = header

            selectDiv.append(textInput)
            selectDiv.append(select)
        }


        // <input type="text"/>
        // <select name="x" id="x-select">
        //     <option value=""></option>
        // </select>
        selectDiv

        return r.value
    })
})