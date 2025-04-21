(() => {
    const deepFreeze = (obj) => {
        if (obj === null || typeof obj !== "object") {
            return obj;
        }
        Object.keys(obj).forEach((key) => {
            deepFreeze(obj[key]);
        });
        return Object.freeze(obj);
    };
    async function dialUp() {
        ws = new WebSocket("/dispatch");
        prom = new Promise((resolve, reject) => {
            ws.onclose = (ev) => { reject(ev); }
            ws.onopen = (ev) => { resolve(ws); }
        });
        return await prom;
    }
    async function delay(timeout) {
        return new Promise((resolve, reject) => {
            setTimeout(resolve, timeout);
        })
    }
    let globalSocket = null;
    let pendingBuffer = [];
    let handlers = new Map();
    async function maintainConnection() {
        while (true) {
            try {
                let ws = await dialUp();
                ws.onmessage = (ev) => {
                    let payload = deepFreeze(JSON.parse(ev.data));
                    let subs = handlers.get(payload.meta.kind);
                    if (subs) {
                        for (s of subs) {
                            if (s !== null) {
                                s(payload);
                            }
                        }
                    }
                }
                console.info('connected', ws);
                let closed = new Promise((resolve, reject) => {
                    ws.onclose = (ev) => { resolve(ev); }
                });
                pendingBuffer.forEach(msg => {
                    ws.send(msg);
                });
                globalSocket = ws;
                pendingBuffer = [];
                await closed;
                globalSocket = null;
                console.info('disconnected', ws);
            } catch {
                await delay(1000);
            }
        }
    }
    window.GlobalSocket = Object.freeze({
        push: function(value) {
            if (globalSocket) {
                globalSocket.send(JSON.stringify(value));
            } else if (pendingBuffer.length < 1000) {
                pendingBuffer.push(JSON.stringify(value));
            }
        },
        next: function(kind, timeout) {
            let subs = handlers.get(kind);
            if (!subs) {
                subs = [];
                handlers.set(kind, subs);
            }
            let prom = new Promise((resolve, reject) => {
                let idx = -1;
                let fn = (val) => {
                    resolve(val);
                    subs[idx] = null;
                };
                if (timeout > 0) {
                    let fnTimeout = async () => {
                        await delay(timeout);
                        reject('timeout');
                        subs[idx] = null;
                    }
                    fnTimeout(); // trigger timeout in background
                }
                idx = subs.findIndex(v => v === null);
                if (idx === -1) {
                    idx = subs.push(fn) - 1;
                } else {
                    subs[idx] = fn;
                }
            });
            return prom;
        }
    });

    let _ = maintainConnection();
})()