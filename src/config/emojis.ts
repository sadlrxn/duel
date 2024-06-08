export interface DuelEmoji {
  id: string;
  type: 'gif' | 'png' | 'jfif';
}

export const duelEmojis: DuelEmoji[] = [
  {
    id: 'pepe_angry_sword_left',
    type: 'gif'
  },
  {
    id: 'pepe_angry_sword_right',
    type: 'gif'
  },
  {
    id: 'pepe_dance_happy',
    type: 'gif'
  },
  {
    id: 'pepe_hands',
    type: 'png'
  },
  {
    id: 'pepe_hmmm',
    type: 'png'
  },
  {
    id: 'pepe_money',
    type: 'gif'
  },
  {
    id: 'pepe_ogiggles',
    type: 'gif'
  },
  {
    id: 'pepe_saber_1',
    type: 'png'
  },
  {
    id: 'pepe_saber_2',
    type: 'png'
  },
  {
    id: 'pepe_shoot',
    type: 'gif'
  },
  {
    id: 'pepe_smirk',
    type: 'png'
  },
  {
    id: 'pepe_wait',
    type: 'png'
  },
  {
    id: 'pepe_pray',
    type: 'png'
  },
  {
    id: 'pepe_hype',
    type: 'gif'
  },
  {
    id: 'pepe_clap',
    type: 'gif'
  },
  {
    id: 'pepe_face_palm',
    type: 'png'
  },
  {
    id: 'pepe_dab',
    type: 'gif'
  },
  {
    id: 'pepe_irritated',
    type: 'gif'
  },
  {
    id: 'pepe_love',
    type: 'png'
  },
  {
    id: 'kekw',
    type: 'png'
  },
  {
    id: 'pepe_gamble',
    type: 'gif'
  },
  {
    id: 'pepe_sniper',
    type: 'png'
  },
  {
    id: 'pepe_yikes',
    type: 'png'
  },
  {
    id: 'pepe_santa',
    type: 'png'
  },
  {
    id: 'pepe_cool',
    type: 'png'
  },
  {
    id: 'pepe_blink',
    type: 'jfif'
  },
  {
    id: 'pepe_cool2',
    type: 'png'
  },
  {
    id: 'pepe_rain',
    type: 'gif'
  },
  {
    id: 'pepe_mcd',
    type: 'png'
  },
  {
    id: 'pepe_cheers',
    type: 'gif'
  },
  {
    id: 'pepe_poncorn',
    type: 'gif'
  },
  {
    id: 'pepe_business',
    type: 'png'
  },
  {
    id: 'pepe_creditcard',
    type: 'gif'
  },
  {
    id: 'pepe_L',
    type: 'png'
  },
  {
    id: 'pepe_clown',
    type: 'png'
  },
  {
    id: 'chad',
    type: 'png'
  },
  {
    id: 'pepe_juice',
    type: 'png'
  },
  {
    id: 'pepe_fishing',
    type: 'png'
  },
  {
    id: 'pepe_lfg',
    type: 'gif'
  },
  {
    id: 'pepe_hug',
    type: 'png'
  },
  {
    id: 'pepe_rocket',
    type: 'gif'
  },
  {
    id: 'astronaut',
    type: 'png'
  }
];

export const generateDuelEmojiUrl = (emoji: DuelEmoji): string => {
  return `https://duelana-bucket.s3.us-east-2.amazonaws.com/emoji/${emoji.id}.${emoji.type}`;
};
