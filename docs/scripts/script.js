function goos() {
    const platform = navigator.platform;
    if (!platform) {
        return '';
    }
    if (platform.indexOf('Win') !== -1) {
        return 'windows';
    }
    if (platform.indexOf('Mac') !== -1) {
        return 'darwin';
    }
    if (platform.indexOf('Linux') !== -1) {
        return 'linux';
    }
    if (platform.indexOf('FreeBSD') !== -1) {
        return 'freebsd';
    }
    if (platform.indexOf('OpenBSD') !== -1) {
        return 'openbsd';
    }
    if (platform.indexOf('SunOS') !== -1) {
        return 'solaris';
    }
    if (platform.indexOf('Android') !== -1) {
        return 'android';
    }
    return '';
}

function defaultLanguage() {
    const threeToTwo = {
        'eng': 'en',
        'jpn': 'ja',
    };
    const lang = navigator.language.match(/^[a-zA-Z]{2,3}/)[0];
    if (lang.length === 3) {
        lang = threeToTwo[lang];
    }
    if (new Set(['en', 'ja']).has(lang)) {
        return lang;
    }
    return 'en';
}

function currentLanguage() {
    const lang = localStorage.getItem('language');
    if (lang) {
        return lang;
    }
    return defaultLanguage();
}

function languageName(code) {
    return {
        'en': 'English',
        'ja': '日本語',
    }[code];
}

// 'darwin' vs 'darwin' => true
// 'darwin' vs 'darwin,amd64' => true
// 'darwin,amd64' vs 'darwin,amd64' => true
// 'darwin,amd64' vs 'darwin' => false
// 'darwin,!amd64' vs 'darwin' => true
function matchesTags(tags, given) {
    const givenSet = new Set(given.split(','));
    loopTerm:
    for (const term of tags.split(' ')) {
        for (const q of term.split(',')) {
            if (q === '') {
                continue;
            }
            if (q.startsWith('!')) {
                if (givenSet.has(q.substring(1)))  {
                    continue loopTerm;
                }
            } else {
                if (!givenSet.has(q)) {
                    continue loopTerm;
                }
            }
        }
        return true;
    }
    return false;
}

function updateCode() {
    for (const e of document.querySelectorAll('pre')) {
        if (!e.dataset['codesrc']) {
            for (const code of e.querySelectorAll('code')) {
                addCommentStyle(code);
            }
            continue;
        }
        (e => {
            fetch(e.dataset['codesrc']).then(r => {
                return r.text();
            }).then(text => {
                if (e.dataset['codelinerange']) {
                    const m = e.dataset['codelinerange'].match(/^(\d+)(-(\d+))?$/);
                    start = parseInt(m[1], 10) - 1;
                    end = start;
                    if (m.length >= 4 && m[3] !== undefined) {
                        end = parseInt(m[3], 10) - 1;
                    }
                    const lines = text.split("\n");
                    text = lines.slice(start, end + 1).join("\n");
                }

                const code = document.createElement('code');
                if (!e.dataset['codelinerange']) {
                    text = text.trim();
                }
                code.textContent = text;
                addCommentStyle(code);
                e.appendChild(code);
            });
        })(e);
    }
}

function addCommentStyle(code) {
    if (code.childNodes.length !== 1) {
        return;
    }
    const text = code.childNodes[0];
    if (text.nodeType !== Node.TEXT_NODE) {
        return;
    }
    code.textContent = '';
    for (const line of text.wholeText.split('\n')) {
        if (!/^\s*\/\//.test(line)) {
            code.appendChild(document.createTextNode(line + '\n'))
            continue;
        }
        const span = document.createElement('span');
        span.classList.add('comment');
        span.textContent = line + '\n'
        code.appendChild(span);
    }
}

function updateImages() {
    for (const e of document.querySelectorAll('p.img')) {
        const img = e.querySelector('img, iframe, video');
        if (!img.complete) {
            const f = () => {
                adjustHeight(img);
                img.removeEventListener('load', f);
            };
            img.addEventListener('load', f);
            continue;
        }
        adjustHeight(img);
    }
}

function adjustHeight(e) {
    // For small diplays, shrink the iframe with keeping its aspect ratio.
    if (e.tagName === 'IFRAME') {
        if (e.clientWidth < e.width) {
            const width = e.clientWidth;
            const ratio = e.height / e.width;
            const height = Math.ceil(width * ratio);
            e.style.height = `${height}px`;
        }
    }

    const unit = 24;
    const height = ~~(((e.clientHeight-1) / unit) + 1) * unit;
    e.parentNode.style.height = `${height}px`;
}

let tocLevel = 4;

function disableTOC() {
    tocLevel = -1;
}

function setTOCLevel(n) {
    tocLevel = n;
}

function updateTOC() {
    let toc = document.querySelector('.toc');
    if (toc !== null) {
        toc.parentNode.removeChild(toc);
    }

    let query = [];
    for (let l = 2; l <= tocLevel; l++) {
        query.push(`article h${l}`);
    }
    if (query.length === 0) {
        return;
    }

    let headers = document.querySelectorAll(query.join(', '));
    for (const header of headers) {
        // https://www.w3.org/TR/html51/dom.html#the-id-attribute
        // The value must be unique amongst all the IDs in the element’s home subtree and must contain at least one
        // character. The value must not contain any space characters.
        header.id = header.textContent.replace(/\s/mg, '_');
    }
    headers = Array.prototype.filter.call(headers, e => {
        return e.offsetParent !== null;
    });
    if (headers.length === 0) {
        return;
    }

    // Create TOC tree.
    toc = document.createElement('div');
    toc.classList.add('toc');
    toc.classList.add('grid-container');
    const gridItem = document.createElement('div');

    gridItem.classList.add('grid-item-4');
    toc.appendChild(gridItem);

    const ul = document.createElement('ul');
    gridItem.appendChild(ul);
    const stack = [ul];

    let last = null;
    for (const header of headers) {
        if (last && last.tagName !== header.tagName) {
            const diff = parseInt(last.tagName.substring(1), 10) - parseInt(header.tagName.substring(1), 10);
            if (diff < 0) {
                const ul = document.createElement('ul');
                const lis = stack[stack.length - 1].querySelectorAll('li');
                lis[lis.length - 1].appendChild(ul);
                stack.push(ul);
            } else {
                for (let i = 0; i < diff; i++) {
                    stack.pop();
                }
            }
        }
        const li = document.createElement('li');
        const a = document.createElement('a');
        a.textContent = header.textContent;
        a.href = `#${header.id}`;
        li.appendChild(a);
        stack[stack.length - 1].appendChild(li);
        last = header;
    }

    const h2s = document.querySelectorAll('main h2');
    for (const h2 of h2s) {
        if (h2.offsetParent === null) {
            continue
        }
        h2.parentNode.insertBefore(toc, h2);
        return;
    }
}

function updateBody() {
    const input = document.querySelector('input#sidemenu');
    // input is null e.g. on the 404 page.
    if (input === null) {
        return;
    }
    if (input.checked) {
        document.body.style.overflow = 'hidden';
    } else {
        document.body.style.overflow = 'visible';
    }
}

function updateCSS() {
    // Trick to override vh unit for mobile platforms.
    // See https://css-tricks.com/the-trick-to-viewport-units-on-mobile/
    const vh = window.innerHeight * 0.01;
    document.documentElement.style.setProperty('--vh', `${vh}px`);
}

function initLanguageSelector() {
    const selector = document.querySelector('#language');
    if (!selector) {
        return;
    }
    const langs = new Set();
    for (const e of document.querySelectorAll('article > *[lang]')) {
        langs.add(e.lang);
    }
    if (langs.size) {
        const sortedLangs = Array.from(langs).sort((a, b) => {
            if (a === 'en') {
                return -1;
            }
            if (b === 'en') {
                return 1;
            }
            if (a < b) {
                return -1;
            }
            if (a > b) {
                return 1;
            }
            return 0;
        });
        for (const lang of sortedLangs) {
            const a = document.createElement('a');
            a.href = '#!';
            a.dataset.lang = lang;
            a.addEventListener('click', e => {
                updateLanguage(lang);
                localStorage.setItem('language', lang);
                e.preventDefault();
            });
            a.textContent = languageName(lang);
            const li = document.createElement('li');
            li.appendChild(a);
            selector.appendChild(li);
        }
    } else {
        selector.style.display = 'none';
    }
}

function updateLanguage(lang) {
    for (const e of document.querySelectorAll('*[lang]')) {
        if (e.lang === lang) {
            e.classList.remove('hiddenlang');
        } else {
            e.classList.add('hiddenlang');
        }
    }
    for (const e of document.querySelectorAll('#language a')) {
        if (e.dataset.lang === lang) {
            e.classList.add('active');
        } else {
            e.classList.remove('active');
        }
    }
}

window.addEventListener('DOMContentLoaded', () => {
    initLanguageSelector();

    updateCode();
    updateImages();
    updateBody();
    updateCSS();
    updateLanguage(currentLanguage());
    updateTOC();

    const sidemenu = document.querySelector('input#sidemenu');
    if (sidemenu !== null) {
        sidemenu.addEventListener('change', updateBody);
    }

    if (typeof katex !== 'undefined') {
        for (const e of document.querySelectorAll('p.math')) {
            const div = document.createElement('div');
            const text = e.textContent;
            e.textContent = '';
            e.appendChild(div);
            katex.render(text, div, {
                displayMode: true,
                strict: true,
            });
            adjustHeight(div);
        }
        for (const e of document.querySelectorAll('span.math')) {
            katex.render(e.textContent, e, {
                displayMode: false,
                strict: true,
            });
        }
    }

    // Twitter
    // https://developer.twitter.com/en/docs/twitter-for-websites/javascript-api/guides/set-up-twitter-for-websites
    if (document.querySelectorAll('blockquote.twitter-tweet').length > 0) {
        window.twttr = ((d, s, id) => {
            var js, fjs = d.getElementsByTagName(s)[0], t = window.twttr || {};
            if (d.getElementById(id)) {
                return t;
            }
            js = d.createElement(s);
            js.id = id;
            js.src = 'https://platform.twitter.com/widgets.js';
            fjs.parentNode.insertBefore(js, fjs);
            t._e = [];
            t.ready = f => {
                t._e.push(f);
            };
            return t;
        })(document, 'script', 'twitter-wjs');
        twttr.ready(() => {
            twttr.events.bind('loaded', e => {
                for (const w of e.widgets) {
                    adjustHeight(w);
                }
            });
        });
    }
});

window.addEventListener('resize', () => {
    updateImages();
    updateCSS();
});
