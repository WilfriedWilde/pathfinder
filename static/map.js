const svgObject = document.getElementById("map");

svgObject.addEventListener("load", () => {
    const svg = svgObject.contentDocument.querySelector("svg");
    const svgWindow = svg.ownerDocument.defaultView;
    const viewBox = svg.viewBox.baseVal;

    const prevButton = document.getElementById("prev");
    const nextButton = document.getElementById("next");
    const turnLabel = document.getElementById("turnLabel");

    const TRAIN_X = 10;
    const TRAIN_Y = 10;
    const LABEL_X = 17;
    const LABEL_Y = 27;
    const MOVE_DURATION = 375;

    let dragging = false;
    let lastX = 0;
    let lastY = 0;
    let currentTurn = 0;

    svg.addEventListener("wheel", e => {
        e.preventDefault();

        const rect = svg.getBoundingClientRect();

        const mouseX = (e.clientX - rect.left) / rect.width;
        const mouseY = (e.clientY - rect.top) / rect.height;

        const oldWidth = viewBox.width;
        const oldHeight = viewBox.height;

        if (e.deltaY < 0) {
            viewBox.width /= 1.2;
            viewBox.height /= 1.2;
        } else {
            viewBox.width *= 1.2;
            viewBox.height *= 1.2;
        }

        viewBox.x += (oldWidth - viewBox.width) * mouseX;
        viewBox.y += (oldHeight - viewBox.height) * mouseY;
    }, { passive: false });

    svg.addEventListener("mousedown", e => {
        e.preventDefault();
        dragging = true;
        lastX = e.clientX;
        lastY = e.clientY;
    });

    svgWindow.addEventListener("mouseup", () => {
        dragging = false;
    });

    svgWindow.addEventListener("mousemove", e => {
        if (!dragging) return;

        const dx = e.clientX - lastX;
        const dy = e.clientY - lastY;
        const svgRect = svg.getBoundingClientRect();

        viewBox.x -= dx * (viewBox.width / svgRect.width);
        viewBox.y -= dy * (viewBox.height / svgRect.height);

        lastX = e.clientX;
        lastY = e.clientY;
    });

    function updateButtons() {
        turnLabel.textContent = `Turn ${currentTurn} / ${turns.length - 1}`;

        prevButton.disabled = currentTurn === 0;
        nextButton.disabled = currentTurn === turns.length - 1;

        prevButton.style.color = currentTurn === 0 ? "gray" : "#0A0A0A";
        nextButton.style.color = currentTurn === turns.length - 1 ? "gray" : "#0A0A0A";
    }

    let trainAnimations = {};
    let trainState = {};


    function setTrainPosition(train, x, y) {
        train.querySelector("circle").setAttribute("cx", x + TRAIN_X);
        train.querySelector("circle").setAttribute("cy", y + TRAIN_Y);
        train.querySelector("text").setAttribute("x", x + LABEL_X);
        train.querySelector("text").setAttribute("y", y + LABEL_Y);
    }

    function animateTrain(trainName, train, fromX, fromY, toX, toY) {
    if (trainAnimations[trainName]) {
        cancelAnimationFrame(trainAnimations[trainName]);
    }

    const start = performance.now();

    function step(now) {
        const t = Math.min((now - start) / MOVE_DURATION, 1);

        const x = fromX + (toX - fromX) * t;
        const y = fromY + (toY - fromY) * t;
        setTrainPosition(train, x, y);

        if (t < 1) {
            trainAnimations[trainName] = requestAnimationFrame(step);
        } else {
            delete trainAnimations[trainName];
        }
    }

    trainAnimations[trainName] = requestAnimationFrame(step);
    }

    function showTurn(turnIndex) {
    const activeTrains = new Set();

    if (turnIndex >= 0) {
        const entries = turns[turnIndex].split(" ");

        for (const entry of entries) {
            const [trainName, stationName] = entry.split("-");
            activeTrains.add(trainName);

            const station = svg.getElementById(stationName);
            const train = svg.getElementById(trainName);

            const targetX = station.cx.baseVal.value;
            const targetY = station.cy.baseVal.value;

            const prevState = trainState[trainName];
            train.setAttribute("visibility", "visible");

            if (!prevState || !prevState.visible) {
                setTrainPosition(train, targetX, targetY);
            } else if (prevState.x !== targetX || prevState.y !== targetY) {
                animateTrain(trainName, train, prevState.x, prevState.y, targetX, targetY);
            }

            trainState[trainName] = { x: targetX, y: targetY, visible: true };
        }
    }

    // hide any train not present this turn
    svg.querySelectorAll(".train").forEach(train => {
        const trainName = train.id;
        if (!activeTrains.has(trainName)) {
            train.setAttribute("visibility", "hidden");

            if (trainAnimations[trainName]) {
                cancelAnimationFrame(trainAnimations[trainName]);
                delete trainAnimations[trainName];
            }
            if (trainState[trainName]) {
                trainState[trainName].visible = false;
            }
        }
    });
}

    prevButton.onclick = () => {
        if (currentTurn === 0) return;

        currentTurn--;
        showTurn(currentTurn);
        updateButtons();
    };

    nextButton.onclick = () => {
        if (currentTurn === turns.length - 1) return;

        currentTurn++;
        showTurn(currentTurn);
        updateButtons();
    };

    window.addEventListener("keydown", e => {
        if (e.key === "ArrowLeft") {
            e.preventDefault();
            prevButton.click();
        } else if (e.key === "ArrowRight") {
            e.preventDefault();
            nextButton.click();
        }
    });

    showTurn(currentTurn);
    updateButtons();
});
