var codeinput = {
    autocompleteMutex: false,
    editor : undefined,
    init : function (){
        editors  = $("#editor")
        codeinput.editor = CodeMirror.fromTextArea(document.getElementById("editor"), {

            matchBrackets: true,
            indentUnit: 8,
            tabSize: 8,
            indentWithTabs: true,
styleActiveLine: true,
            mode: "text/x-go",
            lineNumbers: true,  //显示行号
            extraKeys: {
                "Ctrl-[": "autocompleteAnyWord",
                ".": "autocompleteAfterDot"
            }
        });

        CodeMirror.registerHelper("hint", "go", function (editor) {
            var word = /[\w$]+/;

            var cur = editor.getCursor(), curLine = editor.getLine(cur.line);

            var start = cur.ch, end = start;
            while (end < curLine.length && word.test(curLine.charAt(end))) {
                ++end;
            }
            while (start && word.test(curLine.charAt(start - 1))) {
                --start;
            }

            code = editor.getValue();
            cursorLine = cur.line;
            cursorCh = cur.ch;

            var autocompleteHints = [];

            if (codeinput.autocompleteMutex && editor.state.completionActive) {
                return;
            }

            codeinput.autocompleteMutex = true;
            $.ajax({
                async: false, // 同步执行
                type: 'POST',
                url: '/goimv/goenv/autocomplete',
                data: {code: code,cursorCh:cursorCh,cursorLine:cursorLine},
                dataType: "json",
                success:   function (data) {
                    var autocompleteArray = data[1];

                    if (autocompleteArray) {
                        for (var i = 0; i < autocompleteArray.length; i++) {
                            var displayText = '',
                                text = autocompleteArray[i].name;

                            switch (autocompleteArray[i].class) {
                                case "type":
                                    displayText = // + autocompleteArray[i].class
                                           autocompleteArray[i].name + ''
                                        + autocompleteArray[i].type + '';
                                    break;
                                case "const":
                                    displayText = // + autocompleteArray[i].class
                                         '' + autocompleteArray[i].name + '    '
                                        + autocompleteArray[i].type + '';
                                    break;
                                case "var":
                                    displayText = // + autocompleteArray[i].class
                                         '' + autocompleteArray[i].name + '    '
                                        + autocompleteArray[i].type + '';
                                    break;
                                case "package":
                                    displayText = // + autocompleteArray[i].class
                                         '' + autocompleteArray[i].name + ''
                                        + autocompleteArray[i].type + '';
                                    break;
                                case "func":
                                    displayText = // + autocompleteArray[i].class
                                         '' + autocompleteArray[i].name + ''
                                        + autocompleteArray[i].type.substring(4) + '';
                                    text += '()';
                                    break;
                                default:
                                    console.warn("Can't handle autocomplete [" + autocompleteArray[i].class + "]");
                                    break;
                            }

                            autocompleteHints[i] = {
                                displayText: displayText,
                                text: text
                            };
                        }
                    }
                }});

            setTimeout(function () {
               codeinput.autocompleteMutex = false;
            }, 20);

            return {list: autocompleteHints, from: CodeMirror.Pos(cur.line, start), to: CodeMirror.Pos(cur.line, end)};
        });

        CodeMirror.commands.autocompleteAnyWord = function (cm) {
            cm.showHint({hint: CodeMirror.hint.auto});
        };

        codeinput.editor.on('changes', function (cm) {
            $("#url").html("");
        });


        CodeMirror.commands.autocompleteAfterDot = function (cm) {
            var mode = cm.getMode();
            if (mode && "go" !== mode.name) {
                return CodeMirror.Pass;
            }

            var token = cm.getTokenAt(cm.getCursor());

            if ("comment" === token.type || "string" === token.type) {
                return CodeMirror.Pass;
            }

            setTimeout(function () {
                if (!cm.state.completionActive) {
                    cm.showHint({hint: CodeMirror.hint.go, completeSingle: false});
                }
            }, 50);

            return CodeMirror.Pass;
        };

        $.post('/goimv/problem/list',{},function(data){
            json = JSON.parse(data)
            problems = json['Data']
            for(i =0; i < problems.length; i++ ){
                $('.problem-list').append('<li class="layui-nav-item"><a id=list-item'+i+' class="list-group-item" title='+problems[i]+'>'+problems[i]+'</ a></li>')
                $("#list-item"+i).on("click",{title:problems[i]},function(event){
                    title = event.data.title
                    $.post("/goimv/problem/content",{title:title},function(data){
                        json = JSON.parse(data)
                        content = json["Data"]['content']
                        $(".content").html(marked(content))
                        codeinput.editor.setValue(json['Data']['template'])
                        $(".editor-textarea").attr("title",title)
                    })
                })
            }
        })

    },
    run : function() {
        code = codeinput.editor.getValue()
        title = $(".editor-textarea").attr('title')
        $(".result").html("正在连接服务器......")
        $.post("/goimv/goenv/save",{code:code,title:title},function(data){
            json = JSON.parse(data)
            codeinput.editor.setValue(json['Data']['code'])
            if(json['Errno'] == 0){

                $.post("/goimv/goenv/run",{code:code,title:title,cmd:"run"},function(data){
                    json = JSON.parse(data)
                    $(".result").html(json['Data'])
                })
            }
        })
    },
    test : function() {
        code = codeinput.editor.getValue()
        title = $(".editor-textarea").attr('title')
        $(".result").html("正在连接服务器......")
        $.post("/goimv/goenv/save",{code:code,title:title},function(data){
            json = JSON.parse(data)
            codeinput.editor.setValue(json['Data']['code'])
            if(json['Errno'] == 0){
                $.post("/goimv/goenv/run",{code:code,title:title,cmd:"test"},function(data){
                    json = JSON.parse(data)
                    $(".result").html(json['Data'])
                })
            }
        })

    },


    save :function() {
        code = codeinput.editor.getValue()
        title = $(".editor-textarea").attr('title')
        $.post("/goimv/goenv/save",{code:code,title:title},function(data){
            json = JSON.parse(data)
            codeinput.editor.setValue(json['Data']['code'])

        })
    }
}

$(document).ready(function () {
    codeinput.init();
});
