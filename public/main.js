let options = new Map()

options.set("integer",{
        "upper": 0,
        "lower": 0
    })
options.set("decimal", {
        "digits": 5,
        "commapos": 2
})
options.set("varchar", {
    "length": 50,
    "regex": ""
})

function updateTypeParams(header){
    let fieldInputDiv = document.getElementById(header)
    let columnType = document.getElementById("columnType_"+header)
    let columnTypeParamsDiv = document.getElementById("columnTypeParamsDiv_"+header)
    console.log("columnTypeParamsDiv_"+header)
    columnTypeParamsDiv.innerHTML = ""
    console.log(columnType)
    switch(columnType.value){
        case "varchar":{
            let varcharLengthInput = document.createElement("input")
            varcharLengthInput.type = "number"
            varcharLengthInput.name = "varcharLength_"+header
            varcharLengthInput.id= "varcharLength_"+header
            columnTypeParamsDiv.append(varcharLengthInput)
            break;
        }
        case "decimal":{
            let amountOfDigits = document.createElement("input")
            amountOfDigits.type = "number"
            amountOfDigits.name = "amountOfDigits_"+header
            amountOfDigits.id = "amountOfDigits_"+header
            columnTypeParamsDiv.append(amountOfDigits)
            let commaPosition= document.createElement("input")
            commaPosition.type = "number"
            commaPosition.name = "amountOfDigits_"+header
            commaPosition.id = "amountOfDigits_"+header
            columnTypeParamsDiv.append(commaPosition)
            break
        }
        case "integer":{
            let lowerBound = document.createElement("input")
            lowerBound.type = "number"
            lowerBound.name = "lowerBound_"+header
            lowerBound.id = "lowerBound_"+header
            columnTypeParamsDiv.append(lowerBound)
            let upperBound= document.createElement("input")
            upperBound.type = "number"
            upperBound.name = "lowerBound_"+header
            upperBound.id = "lowerBound_"+header
            columnTypeParamsDiv.append(upperBound)
        }
    }
    return
}


let fileInput = document.getElementById("fileInput")
let selectDiv = document.getElementById("selectDiv")
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
        
        let fieldInputList = document.getElementById("selectDiv")

        for(let i = 0; i < headers.length; i++){
            let header = headers[i].trim()
            let fieldInputDiv = document.createElement("div")
            fieldInputDiv.id = header

            let columnNameInput = document.createElement("input")
            columnNameInput.name= "columnName_" + header
            columnNameInput.id= "columnName_" + header
            columnNameInput.value = header
            fieldInputDiv.append(columnNameInput)


            let columnTypeInput = document.createElement("select")
            columnTypeInput.name = "columnType_" + header
            columnTypeInput.id= "columnType_" + header
            for(let [key, value] of options){
                let option = document.createElement("option")
                option.value = key
                option.text = key
                columnTypeInput.append(option)
                option.selected = true
            }
            console.log(columnTypeInput.value)

            fieldInputDiv.append(columnTypeInput)
            

            let columnTypeParamsDiv = document.createElement("div")
            columnTypeParamsDiv.id = "columnTypeParamsDiv_" + header
            fieldInputDiv.append(columnTypeParamsDiv)



            fieldInputList.append(fieldInputDiv)
            updateTypeParams(header)
            columnTypeInput.addEventListener("change", (e)=>updateTypeParams(header))
        }
        


        // <input type="text"/>
        // <select name="x" id="x-select">
        //     <option value=""></option>
        // </select>

        return r.value
    })
})