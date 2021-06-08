"use strict";
var __awaiter = (this && this.__awaiter) || function (thisArg, _arguments, P, generator) {
    function adopt(value) { return value instanceof P ? value : new P(function (resolve) { resolve(value); }); }
    return new (P || (P = Promise))(function (resolve, reject) {
        function fulfilled(value) { try { step(generator.next(value)); } catch (e) { reject(e); } }
        function rejected(value) { try { step(generator["throw"](value)); } catch (e) { reject(e); } }
        function step(result) { result.done ? resolve(result.value) : adopt(result.value).then(fulfilled, rejected); }
        step((generator = generator.apply(thisArg, _arguments || [])).next());
    });
};
var __importDefault = (this && this.__importDefault) || function (mod) {
    return (mod && mod.__esModule) ? mod : { "default": mod };
};
Object.defineProperty(exports, "__esModule", { value: true });
exports.createTendermintSocketPoll = exports.TendermintSocketPoll = void 0;
const eventemitter2_1 = require("eventemitter2");
const axios_1 = __importDefault(require("axios"));
function fetchBlock(url) {
    return __awaiter(this, void 0, void 0, function* () {
        const res = yield axios_1.default.get(url);
        return res.data;
    });
}
function TendermintSocketPoll({ apiUrl, fetcher = fetchBlock, pollInterval = 5000, }) {
    const emitter = new eventemitter2_1.EventEmitter2();
    function pollBlock(height) {
        return __awaiter(this, void 0, void 0, function* () {
            const query = typeof height !== "undefined" ? `?height=${height}` : "";
            return yield fetcher(`${apiUrl.replace(/\/$/, "")}/block${query}`);
        });
    }
    // Process a block and emit events based on that block
    function processData(blockData) {
        emitter.emit("NewBlock", blockData);
        const txs = blockData.result.block.data.txs;
        if (txs) {
            txs.forEach((tx) => {
                // TODO: Not sure if we should/can add more tx information here as all we have is the encoded tx - can we decode it? need to look into it
                emitter.emit("Tx", tx);
            });
        }
        // Return processed blockheight
        return parseInt(blockData.result.block.header.height);
    }
    function sleep(ms) {
        return new Promise((resolve) => setTimeout(resolve, ms));
    }
    /**
     * Get all blocks that havent been processed since last known blockheight
     * @param height last processed blockheight or null for no new blockheight
     */
    function getBlocksToProcess(height) {
        return __awaiter(this, void 0, void 0, function* () {
            const blockData = yield pollBlock();
            const newHeight = parseInt(blockData.result.block.header.height);
            // If height is null this is the first poll so process this block
            if (height === null) {
                return [blockData];
            }
            // no new data don't process any blocks
            if (newHeight === height) {
                return [];
            }
            // There are blocks to be processed build up a list of them
            let heightDiff = newHeight - height;
            const blocks = [blockData];
            for (let i = heightDiff - 1; i > 0; --i) {
                const interimData = yield pollBlock(height + i);
                blocks.unshift(interimData);
            }
            return blocks;
        });
    }
    let polling = false;
    function startPoll() {
        return __awaiter(this, void 0, void 0, function* () {
            // If already polling dont poll again
            if (polling)
                return;
            polling = true;
            let height = null;
            // Loop while we are polling
            while (polling) {
                // First we get a list of blocks to process
                const blocks = yield getBlocksToProcess(height);
                // Then we process them updating the height
                for (let block of blocks) {
                    height = processData(block);
                }
                // Then let's wait for a poll interval
                yield sleep(pollInterval);
            }
        });
    }
    function stopPoll() {
        polling = false;
    }
    return {
        on(event, handler) {
            if (!emitter.hasListeners()) {
                startPoll();
            }
            emitter.on(event, handler);
        },
        off(event, handler) {
            emitter.off(event, handler);
            if (!emitter.hasListeners()) {
                stopPoll();
            }
        },
    };
}
exports.TendermintSocketPoll = TendermintSocketPoll;
// Make this a singleton to avoid multiple polling
let instance;
function createTendermintSocketPoll(apiUrl) {
    if (!instance) {
        instance = TendermintSocketPoll({ apiUrl });
    }
    return instance;
}
exports.createTendermintSocketPoll = createTendermintSocketPoll;
//# sourceMappingURL=TendermintSocketPoll.js.map