(() => {
    const events = ['click'];
    let dispatchers = new Map();
    events.forEach((v) => {
        dispatchers.set(v, makeDispatcher(v));
    });
    function makeDispatcher(name) {
        return (ev) => {
            let target = ev.target;
            let value = target.getAttribute(`et-${name}`);
            if (!value) {
                return false;
            }
            ev.preventDefault();
            ev.stopPropagation();
            let jvalue = JSON.parse(value);
            let payload = {
                value: jvalue,
                meta: {
                    kind: name,
                    path: computePath(target, []).reverse(),
                }
            }
            if (globalThis.GlobalSocket) {
                globalThis.GlobalSocket.push(payload);
            };
            return false;
        }
    }
    function registerEventHandler(name) {
        let items = document.querySelectorAll(`[et-${name}]`);
        for (let item of items) {
            item.removeEventListener(name, dispatchers.get(name));
            item.addEventListener(name, dispatchers.get(name));
        }
    }
    function registerEventTriggers() {
        for(let evname of events) {
            registerEventHandler(evname);
        }
    }
    function computePath(el, acc) {
        let id = el.getAttribute("id");
        let etID = el.getAttribute("et-id");
        let tag = el.tagName.toLowerCase();
        let obj = {
            tag: tag
        };
        if (id) {
            obj.id = id;
        }
        if (etID) {
            obj.etID = etID;
        }
        acc.push(obj);
        if (el.parentElement && el.parentElement.tagName != 'HTML') {
            return computePath(el.parentElement, acc);
        }
        return acc
    }
    document.addEventListener("DOMContentLoaded", registerEventTriggers);
})()