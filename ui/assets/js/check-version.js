(() => {
    async function delay(timeout) {
        return new Promise((resolve, _) => {
            setTimeout(resolve, timeout);
        });
    }
    let oldVersion;
    async function fetchVersion() {
        while(true) {
            try {
                globalThis.GlobalSocket.push({value: {}, meta: { kind: 'core.app.get-version'}});
                let version = (await globalThis.GlobalSocket.next('core.app.version', 3000)).value;
                if (!oldVersion) {
                    oldVersion = version;
                } else if (oldVersion.version != version.version) {
                    window.location.reload();
                } else {
                    await delay(1000)
                }
            } catch {
                console.error('timeout!');
            }
        }
    }
    fetchVersion();
})()