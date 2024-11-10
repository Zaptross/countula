window.addEventListener("load", function () {
    // A map of rule IDs to rule values
    const ruleData = parseRulesFromSQuery();
    const ruleInputs = Array.from(document.getElementsByTagName("input"));
    initialiseRuleInputs(ruleInputs, ruleData);
    const updateRuleCommand = (ruleData) => {
        const ruleCommandElement = document.getElementById("rule-command");

        let ruleCommand = "/count settings action:set settings:";
        const settings = [];

        const entries = Object.entries(ruleData);
        for (const [ruleId, ruleValue] of entries) {
            settings.push(`${ruleId}:${ruleValue}`);
        }

        ruleCommandElement.innerText = `${ruleCommand} ${settings.join(",")}`;
    };
    const spreadDifferenceAcrossOthersOfType = (
        ruleInputs,
        ruleData,
        ruleType,
        ruleId,
        difference
    ) => {
        const otherRuleIdsOfType = ruleInputs
            .filter((input) => input.dataset.ruleType === ruleType)
            .map((input) => input.dataset.ruleId)
            .filter((id) => id !== ruleId)
            .reduce((acc, id) => {
                if (!acc.includes(id)) {
                    acc.push(id);
                }
                return acc;
            }, []);

        let i = 0;
        while (difference !== 0) {
            if (i >= otherRuleIdsOfType.length) {
                i = 0;
                continue;
            }
            const change = Math.sign(difference);
            const otherRuleId = otherRuleIdsOfType[i];

            // Skip zeroed out rules
            if (ruleData[otherRuleId] <= 0 && change > 0) {
                i++;
                continue;
            }

            ruleData[otherRuleId] -= change;
            difference -= change; // Approach zero
            i++;
        }

        for (const [ruleId, ruleValue] of Object.entries(ruleData)) {
            if (!ruleId || isNaN(ruleValue)) {
                continue;
            }

            const ruleWeightDisplay = document.getElementById(
                "rule-" + ruleId + "-display"
            );
            const ruleRangeDisplay = document.getElementById(
                "rule-" + ruleId + "-range"
            );
            ruleWeightDisplay.value = ruleValue;
            ruleRangeDisplay.value = ruleValue;
        }
    };

    ruleInputs
        .filter((input) => input.type === "range")
        .map((input) => {
            let ruleWeightDisplay;
            {
                const ruleId = input.dataset.ruleId;
                const ruleValue = input.value;
                ruleData[ruleId] = ruleValue;
                ruleWeightDisplay = document.getElementById(
                    "rule-" + ruleId + "-display"
                );
            }
            input.addEventListener("change", (event) => {
                const ruleId = event.target.dataset.ruleId;
                const ruleType = event.target.dataset.ruleType;
                const ruleValue = event.target.value;
                const difference = ruleValue - ruleData[ruleId];
                ruleData[ruleId] = ruleValue;

                ruleWeightDisplay.value = ruleValue;
                spreadDifferenceAcrossOthersOfType(
                    ruleInputs,
                    ruleData,
                    ruleType,
                    ruleId,
                    difference
                );
                updateRuleCommand(ruleData);
            });
        });

    ruleInputs
        .filter((input) => input.type === "number")
        .map((input) => {
            let ruleRangeDisplay = document.getElementById(
                "rule-" + input.dataset.ruleId + "-range"
            );
            input.addEventListener("change", (event) => {
                const ruleId = event.target.dataset.ruleId;
                const ruleType = event.target.dataset.ruleType;
                const ruleValue = event.target.value;
                const difference = ruleValue - ruleData[ruleId];
                ruleData[ruleId] = ruleValue;

                ruleRangeDisplay.value = ruleValue;
                spreadDifferenceAcrossOthersOfType(
                    ruleInputs,
                    ruleData,
                    ruleType,
                    ruleId,
                    difference
                );
                updateRuleCommand(ruleData);
            });
        });

    // Update the rule command on page load
    updateRuleCommand(ruleData);
});

function parseRulesFromSQuery() {
    // Eg; 1:24,2:24,4:20,8:16,16:60,32:56,64:40,128:16,256:28,512:16
    const sQuery = new URLSearchParams(window.location.search).get("s");
    if (!sQuery) {
        return {};
    }

    return sQuery
        .split(",")
        .map((setting) => setting.split(":"))
        .reduce((acc, [id, value]) => ({...acc, [id]: parseInt(value)}), {});
}

// Initialise the rule inputs with the values from the query string
function initialiseRuleInputs(ruleInputs, ruleData) {
    for (const input of ruleInputs) {
        const ruleId = input.dataset.ruleId;
        if (!ruleId) {
            continue;
        }

        const ruleValue = ruleData[ruleId];
        if (ruleValue !== undefined) {
            input.value = ruleValue;
        }
    }
}
