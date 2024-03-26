const playerUrl = "http://localhost:8888/player"
const fetchConfig = {headers: {"Content-Type": "application/json", "Authorization": "somevalue"}};

function $(id) {
    return document.getElementById(id);
}

function createElement(tagName, content) {
    const el = document.createElement(tagName);
    if (typeof content === "string") {
        el.innerHTML = content;
    }
    if (content && typeof content === "object") {
        for (let k in content) {
            el[k] = content[k];
        }
    }
    return el;
}

function withEvent(event, handler) {
    return (el) => {
        el.addEventListener(event, handler);
        return el;
    };
}

function append(target, children) {
    children.forEach(ch => target.appendChild(ch));
}

const state = {
    nameValue: "",
    answer: "",
    players: [],
    guessValue: "",
}

function nameUI() {
    const nameLabel = createElement("label", "Enter Name:");
    const nameInput = withEvent("change", handleNameChange)(createElement("input", {id: "nameInput"}));
    const nameButton = withEvent("click", handleNameSend)(createElement("button", "send"));
    return [nameLabel, nameInput, nameButton];
}

function gamePlayUI() {
    const hints = createElement("div", {id: "hints"});
    const hintsTitle = createElement("h2", "Hints from server");
    append(hints, [hintsTitle]);
    const guessLabel = createElement("label", "Guess the word:")
    const guessInput =  withEvent("change", handleGuessChange)(createElement("input", {id: "guessInput"}));
    const guessButton = withEvent("click", handleGuessSend)(createElement("button", "submit guess"));
    return [hints, guessInput, guessButton];
}

function app() {
    const container = $("app");
    const title = createElement("h1", "Play Word Games");
    const fetcherButton = withEvent("click", getPlayers)(createElement("button", "test fetch"));

    append(container, [title, ...nameUI(), fetcherButton, ...gamePlayUI()]);
}

async function handleNameSend() {
    if (state.nameValue) {
        console.log(`posting value: ${state.nameValue}`);
        await createPlayer();
    }
    state.inputValue = "";
    $("nameInput").value = "";
}

function handleNameChange(ev) {
    state.nameValue = ev.target.value;
}

function handleGuessChange(ev) {
    state.guessValue = ev.target.value;
}

function handleGuessSend() {
    console.log(`posting guess word '${state.guessValue}' to api...`);
    state.guessValue = "";
}

app();

async function getPlayers() {
    const res = await fetch(playerUrl, fetchConfig).then(r => r.json());

    console.log(res);
}

async function createPlayer() {
    const playerRequest = {name: state.nameValue};
    const postConfig = {...fetchConfig, body: JSON.stringify(playerRequest), method: "POST"};
    const res = await fetch(playerUrl, postConfig).then(r => r.text());

    console.log(res);
}