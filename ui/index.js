const playerUrl = "http://localhost:8888/player";
const gameUrl = "http://localhost:8888/game";
const startUrl = "http://localhost:8888/start";
const guessUrl = "http://localhost:8888/guess";
const scoreUrl = "http://localhost:8888/score";
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
    players: {},
    gameId: "",
    guessValue: "",
};

function nameUI() {
    const nameLabel = createElement("label", "Enter Name:");
    const nameInput = withEvent("change", handleNameChange)(createElement("input", {id: "nameInput"}));
    const nameButton = withEvent("click", handleNameSend)(createElement("button", "send player name"));
    return [nameLabel, nameInput, nameButton];
}

function gamePlayUI() {
    const gameButton = withEvent("click", createGame)(createElement("button", "create a game"));
    const startButton = withEvent("click", startGame)(createElement("button", "start game"));
    const hints = createElement("div", {id: "hints"});
    const hintsTitle = createElement("h2", "Hints from server");
    append(hints, [hintsTitle]);

    const guessLabel = createElement("label", "Guess the word:")
    const guessInput =  withEvent("change", handleGuessChange)(createElement("input", {id: "guessInput"}));
    const guessButton = withEvent("click", handleGuessSend)(createElement("button", "submit guess"));
    const scoreboardContainer = createElement("div", {id: "scoreboard"});
    return [gameButton, startButton, hints, guessLabel, guessInput, guessButton, scoreboardContainer];
}

function app() {
    console.log("starting ui...");
    const container = $("app");
    const title = createElement("h1", "Play Word Games");
    const fetcherButton = withEvent("click", getPlayers)(createElement("button", "test fetch players"));

    append(container, [title, ...nameUI(), ...gamePlayUI(), fetcherButton]);
}

async function handleNameSend() {
    if (state.nameValue) {
        console.log(`posting value: ${state.nameValue}`);
        await createPlayer();
    }
    state.nameValue = "";
    $("nameInput").value = "";
}

function handleNameChange(ev) {
    state.nameValue = ev.target.value;
}

function handleGuessChange(ev) {
    state.guessValue = ev.target.value;
}

async function handleGuessSend() {
    try {
        const body = JSON.stringify({guess: state.guessValue, game_id: state.gameId})
        console.log(`posting body ${body} to api...`);
        const postConfig = {
            ...fetchConfig,
            body,
            method: "POST"
        };
        const res = await fetch(guessUrl, postConfig).then(r => r.json());
        makeModal(res.status);
        if (res.status.includes("winner")) {
            $("hints").innerHTML = "<h2>Hints from server</h2>";
            $("guessInput").value = "";
            await getScores();
        }
        console.log(res);
    } catch (e) {
        console.error(e);
    }
}

async function getScores() {
    try {
        const getConfig = {...fetchConfig, method: "GET"};
        const res = await fetch(scoreUrl, getConfig).then(r => r.json());
        $("scoreboard").innerHTML = `<h2>Scoreboard:</h2><pre>${JSON.stringify(res, null, 2)}</pre>`
    } catch (e) {
        console.error(e);
    }
}

async function getPlayers() {
    const res = await fetch(playerUrl, fetchConfig).then(r => r.json());
    console.log(res);
}

async function createPlayer() {
    try {
        const nameValue = state.nameValue
        const playerRequest = {name: nameValue};
        const postConfig = {...fetchConfig, body: JSON.stringify(playerRequest), method: "POST"};
        const res = await fetch(playerUrl, postConfig).then(r => r.json());
        console.log("player created:", res);
        state.players[res.player_id] = nameValue;
        console.log(state);
    } catch (e) {
        console.error(e);
    }
}

async function createGame() {
    try {
        const playerId = Object.keys(state.players).find((k, i) => {
            if (i === 0) {
                return k;
            }
        });
        const postConfig = {...fetchConfig, body: JSON.stringify({player_id: playerId}), method: "POST"};
        const res = await fetch(gameUrl, postConfig).then(r => r.json());
    
        state.gameId = res.game_id;
        console.log("response: ", res);
        console.log("state: ", state);

    } catch (e) {
        console.error(e);
    }
}
async function startGame() {
    try {
        const playerId = Object.keys(state.players).find((k, i) => {
            if (i === 0) {
                return k;
            }
        });
        const postConfig = {
            ...fetchConfig,
            body: JSON.stringify({player_id: playerId, game_id: state.gameId}), method: "POST"
        };
        const res = await fetch(startUrl, postConfig).then(r => r.json());
        console.log(res)
        $("hints").innerHTML = `<h2>Hints from server</h2><pre>${makeHints(res)}</pre>`;
        
    } catch(e) {
        console.error(e);
    }
}

function makeHints({word, meanings}) {
    const wStar = `<p>Word: ${word}</p>`;
    const list = mapMeanings(meanings);
    const ul = `<ul>${list}</ul>`;
    return `<div>${wStar}${ul}</div>`;
  }
  
  function mapMeanings(meanings) {
    return meanings.map((m) => {
      const {partOfSpeech, definitions} = m;
      return `<li><p>Part of speech: ${partOfSpeech}</p>${mapDefinitions(definitions)}`;
    }).join("");
  }
  
  function mapDefinitions(definitions) {
    return definitions.map((d, i) => `<p>${i}. ${d.definition}</p>`).join("");
  }

function removeMe() {
    $("app").removeChild(this.parentNode);
}
function makeModal(msg) {
    const modal = createElement("div", {id: "modal", "style": "border:1px solid #000;position:fixed;top:0;left:0;z-index:2;background-color:#333;color:#fff;padding:1em;"});
    const p = createElement("p", msg);
    const b = withEvent("click", removeMe)(createElement("button", "X"));
    append(modal, [b, p]);
    $("app").appendChild(modal);
}

app();