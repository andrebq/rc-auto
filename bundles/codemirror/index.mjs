import {EditorView, basicSetup} from "codemirror"
import {StreamLanguage} from "@codemirror/language"
import {javascript} from "@codemirror/lang-javascript"
import {python} from "@codemirror/lang-python"
import {go} from "@codemirror/lang-go"
import {lua} from "@codemirror/legacy-modes/mode/lua"
import {shell} from "@codemirror/legacy-modes/mode/shell"
import {oneDark} from "@codemirror/theme-one-dark"

Object.defineProperty(window, 'NewCodeEditor', {
    value: (parent, language) => {
        let langObj = null;
        if (language === 'python') {
            langObj = python()
        } else if (language === 'javascript') {
            langObj = javascript()
        } else if (language === 'shell') {
            langObj = StreamLanguage.define(shell);
        } else if (language === 'lua') {
            langObj = StreamLanguage.define(lua);
        } else if (language === 'go') {
            langObj = go()
        }
        let ext = [basicSetup, oneDark];
        if (langObj) {
            ext.push(langObj);
        }
        let editor = new EditorView({
            extensions: ext,
            parent: parent
        });
        return editor;
    },
    writable: false,
    enumerable: true,
    configurable: false,
})

