import { Tab, TabList } from 'react-tabs';
import { Panel } from 'rc-collapse';
import 'react-tabs/style/react-tabs.css';
import 'rc-collapse/assets/index.css';

import { ReactComponent as CaretIcon } from 'assets/imgs/icons/caret.svg';
import { ReactComponent as WIDIcon } from 'assets/imgs/help/what-is-duelana.svg';
import { ReactComponent as HTPIcon } from 'assets/imgs/help/how-to-play.svg';
import { ReactComponent as UpdateIcon } from 'assets/imgs/help/updates.svg';
import { Text, Box } from 'components';

import {
  StyledCollapse,
  StyledTabs,
  StyledTabPanel,
  StyledLink
} from './styles';

export default function Help() {
  return (
    <Box padding={['30px 12px', '30px 12px', '30px 12px', '30px 25px']}>
      <Text
        fontSize={'38px'}
        fontWeight={700}
        color="#fff"
        textAlign={'center'}
        mb="40px"
      >
        How can we help you ?
      </Text>
      <StyledTabs>
        <TabList>
          <Tab>
            <WIDIcon />
            <Text fontWeight={700} fontSize="20px" mt="10px" color="#fff">
              What is Duel
            </Text>
            <Text fontSize={'14px'} color="#515C6D">
              Learn more about Duel
            </Text>
          </Tab>
          <Tab>
            <HTPIcon />
            <Text fontWeight={700} fontSize="20px" mt="10px" color="#fff">
              How to play
            </Text>
            <Text fontSize={'14px'} color="#515C6D">
              Learn to play our games
            </Text>
          </Tab>
          <Tab>
            <UpdateIcon />
            <Text fontWeight={700} fontSize="20px" mt="10px" color="#fff">
              Updates
            </Text>
            <Text fontSize={'14px'} color="#515C6D">
              Stay up to date with Duel
            </Text>
          </Tab>
        </TabList>

        <StyledTabPanel>
          <StyledCollapse
            accordion={true}
            expandIcon={() => <CaretIcon width={24} height={12} />}
          >
            <Panel header="What is Duel?">
              Duel is a fully licensed online Casino that allows users to wager
              using cryptocurrencies and NFTs as collateral. We are a
              community-driven platform that strives to provide the best online
              gambling experience utilizing the latest innovations in Web3. All
              of our games, features, and UX/UI are made in-house from scratch
              by our dev team.
            </Panel>

            <Panel header="Is there a Fee?">
              Yes, Duel takes a service fee of 2% for Coin Flip and 5% for
              Jackpot.
            </Panel>

            <Panel header="What makes Duel unique?">
              First and foremost, we genuinely care about delivering a great
              gambling experience that you enjoy coming back to time and time
              again. Our mission is to provide users with the most secure,
              legal, and provably-fair system possible. We believe we can
              improve the existing experience by taking advantage of the
              numerous innovations unlocked by a fast, scalable blockchain like
              Solana. In addition to depositing and withdrawing in your favorite
              spl-token, our platform will let you utilize your NFTs in various
              ways as collateral for your wagers. We aim to improve the games we
              all already love with the technologies of the future in a platform
              that encourages the community to engage and bond with each other.
            </Panel>
            <Panel header="Are you launching an NFT collection?">
              Yes, “Duelbots” is a collection of gambling-themed robot duelers;
              there will be 2,223 Duelbots up for grabs. Our Duelbots will offer
              various benefits to holders on our Duel platform, including but
              not limited to revenue sharing, loans, and loot boxes.
            </Panel>
            <Panel header="Where can I find more information?">
              Visit our knowledge base at{' '}
              <StyledLink
                href="https://docs.duel.win"
                rel="noreferrer"
                target={'_blank'}
              >
                https://docs.duel.win
              </StyledLink>{' '}
              and our blog at{' '}
              <StyledLink
                href="https://blog.duel.win"
                rel="noreferrer"
                target={'_blank'}
              >
                https://blog.duel.win
              </StyledLink>{' '}
              for the latest updates.
            </Panel>
            <Panel header="What games are you building?">
              At launch, we will be debuting 3 of our fully functioning and
              in-house developed games; Coin Flip, Jackpot, and Dream Towers.
              Our Jackpot game features our unique NFT wagering system, which
              allows you to use your favorite JPGs! We will continue to expand
              our game offerings based on community feedback, and our goal is to
              evolve into a complete online casino.
            </Panel>
            <Panel header="What is the Rakeback system?">
              The Rakeback system is a loyalty program we have designed to
              reward our members for playing on our Duel platform. You will
              receive an instant rake back of 5% of the house fees you have paid
              while wagering CHIPs on our platform (distributed in the form of
              CHIPs).
            </Panel>
            <Panel header="Is this casino legal?">
              Absolutely. We are a fully licensed casino and have been
              thoroughly audited by the necessary legal entities. Read more
              about why casino licenses are important here:{' '}
              <StyledLink
                href="https://blog.duel.win/articles/why-is-a-casino-license-important/"
                rel="noreferrer"
                target={'_blank'}
              >
                https://blog.duel.win/articles/why-is-a-casino-license-important/
              </StyledLink>
            </Panel>
            <Panel header="How do I know this casino is safe?">
              Rest assured, your funds are completely safe, and our games have
              been vigorously vetted to make sure they are fair. Our casino
              license means that we adhere to strict operating standards and
              keep everything above board 100% of the time.
            </Panel>
            <Panel header="What is 1 CHIP worth?">
              1 CHIP is the equivalent of 1 Dollar USD.
            </Panel>
          </StyledCollapse>
        </StyledTabPanel>

        <StyledTabPanel>
          <StyledCollapse
            accordion={true}
            expandIcon={() => <CaretIcon width={24} height={12} />}
          >
            <Panel header="Coin Flip">
              Like any traditional coinflip, our game consists of Heads(Green)
              and Tails(Purple). To start, you need to pick a side and place
              your wager. Please check out{' '}
              <StyledLink
                href="https://docs.duel.win/games/coin-flip"
                rel="noreferrer"
                target={'_blank'}
              >
                https://docs.duel.win/games/coin-flip
              </StyledLink>{' '}
              for further details.
            </Panel>
            <Panel header="Jackpot">
              Wager both CHIPS or NFTs in one pot to win it all. The more you
              wager, the higher chances of winning. Please check out&nbsp;
              <StyledLink
                href="https://docs.duel.win/games/jackpot"
                rel="noreferrer"
                target={'_blank'}
              >
                https://docs.duel.win/games/jackpot
              </StyledLink>
              &nbsp;for further details.
            </Panel>
            <Panel header="Dream Tower">
              The objective of the game is to select the tile that contains a
              star for each level (9 levels in total). The higher you go, the
              higher your payout multiplier becomes. For a more detailed
              explanation (and visuals) check out our blog:{' '}
              <StyledLink
                href="https://blog.duel.win/updates/dream-tower/"
                rel="noreferrer"
                target={'_blank'}
              >
                https://blog.duel.win/updates/dream-tower/
              </StyledLink>
              .
            </Panel>
          </StyledCollapse>
        </StyledTabPanel>

        <StyledTabPanel>
          <StyledCollapse
            accordion={true}
            expandIcon={() => <CaretIcon width={24} height={12} />}
          >
            <Panel header="Coin Flip Live">
              Duel has launched its first game, Coin Flip.
            </Panel>

            <Panel header="Jackpot Live">
              Duel has launched its second game, Jackpot.
            </Panel>

            <Panel header="Dream Tower Live">
              Duel has launched its third game, Dream Tower.
            </Panel>
          </StyledCollapse>
        </StyledTabPanel>
      </StyledTabs>
    </Box>
  );
}
