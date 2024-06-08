export const neverMatchingRegex = /($a)/;

export const notifyRegex = /(@\{\w+\})/g;
export const duelEmojiRegex = /(:\{\w+\})/g;
export const tipMsgRegex = /(\$ \d+ .*)/g;

// Command
export const tipRegex = /\/tip @?\{?\w+\}?/;
export const detailsRegex = /\/details @?\{?\w+\}?/;
export const muteRegex = /\/mute @?\{?\w+\}?/;
export const unmuteRegex = /\/unmute @?\{?\w+\}?/;
export const banRegex = /\/ban @?\{?\w+\}?/;
export const unbanRegex = /\/unban @?\{?\w+\}?/;
export const setWagerLimitRegex = /\/setWagerLimit \d+/;
export const setMaxLengthRegex = /\/setMaxLength \d+/;
export const rainRegex = /\/rain \d+ \d+/;
