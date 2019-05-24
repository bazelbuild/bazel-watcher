let counter = 0;

export function displayDuration(element) {
    setInterval(() => {
        element.innerText = `${++counter}`;
    }, 1000 /* 1s */);
}
