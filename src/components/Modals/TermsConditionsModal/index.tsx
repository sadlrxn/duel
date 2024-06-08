import React, { useCallback } from "react";

import { ReactComponent as TermsIcon } from "assets/imgs/icons/terms.svg";
import { P, Box, Flex, Modal, ModalProps, Button } from "components";
import { Text } from "components/Text";
import styled from "styled-components";

interface TermsConditionsModalProps extends ModalProps {
  setUserAccepted: React.Dispatch<React.SetStateAction<boolean>>;
  login: () => void;
}

export default function TermsConditionsModal({
  setUserAccepted,
  login,
  ...props
}: TermsConditionsModalProps) {
  const handleAccept = useCallback(() => {
    setUserAccepted(true);
    props.onDismiss?.();
    login();
  }, [login, props, setUserAccepted]);

  return (
    <Modal {...props}>
      <Box
        p={["40px 20px", "40px 30px", "40px 40px", "40px 50px"]}
        background="linear-gradient(180deg, #132031 0%, #1a293c 100%)"
        border="2px solid #7389a9"
        borderRadius="20px"
        maxWidth={997}
      >
        <Container>
          <Flex gap={15} alignItems="center">
            <TermsIcon />
            <Text
              color="#FFFFFF"
              textTransform="uppercase"
              fontWeight={600}
              fontSize={20}
            >
              Terms & Conditions
            </Text>
          </Flex>

          <Box
            overflow="auto"
            pl={20}
            pr={31}
            py={20}
            background="#121C2A"
            borderRadius="11px"
          >
            <Box>
              <Text color="#768BAD" fontWeight={600} fontSize={16}>
                Definitions
              </Text>
              <Text color="#768BAD" fontSize={12}>
                <List>
                  <li>
                    Duel Games Corp. is referred to as &quot;Duel&quot;,
                    &quot;we&quot; or &quot;us&quot;.
                  </li>
                  <li>
                    The Player is referred to as &quot;you&quot; or &quot;the
                    Player&quot;.
                  </li>
                  <li>
                    &quot;Games&quot; means Casino and other games as may from
                    time to time become available on the Websites.
                  </li>
                  <li>
                    &quot;The Website&quot; means https://duel.win through
                    desktop, mobile or other platforms utilized by the Player.
                  </li>
                  <li>1 Solana = 1 SOL = 1000000000 Lamports</li>
                </List>
              </Text>
            </Box>

            <Box my={15}>
              <Text color="#768BAD" fontWeight={600} fontSize={16}>
                1. General
              </Text>
              <Text
                color="#768BAD"
                fontWeight={400}
                fontSize={12}
                lineHeight="19px"
              >
                <P>
                  1.1 | These terms and conditions (&quot;Terms and
                  Conditions&quot;) apply to the usage of games accessible
                  through{" "}
                  <StyledLink
                    href="https://duel.win"
                    rel="noreferrer"
                    target={"_blank"}
                  >
                    https://duel.win
                  </StyledLink>
                </P>
                <P>
                  1.2 | These Terms and Conditions come into force as soon as
                  you complete the registration process, which includes checking
                  the box accepting these Terms and Conditions and successfully
                  creating a member account (&quot;Member Account&quot;). By
                  using any part of the Website following Member Account
                  creation, you agree to these Terms and Conditions applying to
                  the use of the Website.
                </P>
                <P>
                  1.3 | You must read these Terms and Conditions carefully in
                  their entirety before creating a Member Account. If you do not
                  agree with any provision of these Terms and Conditions, you
                  must not create a Member Account or continue to use the
                  Website.
                </P>
                <P>
                  1.4 | We are entitled to make amendments to these Terms and
                  Conditions at any time and without advanced notice. If we make
                  such amendments, we may take appropriate steps to bring such
                  changes to your attention (such as by email or placing a
                  notice on a prominent position on the Website, together with
                  the amended terms and conditions) but it shall be your sole
                  responsibility to check for any amendments, updates and/or
                  modifications. Your continued use of Duel&apos;s services and
                  Website after any such amendment to the Terms and Conditions
                  will be deemed as your acceptance and agreement to be bound by
                  such amendments, updates and/or modifications.
                </P>
                <P>
                  1.5 | These Terms and Conditions may be published in several
                  languages for informational purposes and ease of access by
                  players. The English version is the only legal basis of the
                  relationship between you and us and in the case of any
                  discrepancy with respect to a translation of any kind, the
                  English version of these Terms and Conditions shall prevail.
                </P>
              </Text>
            </Box>

            <Box>
              <Text color="#768BAD" fontWeight={600} fontSize={16}>
                2. Binding Declarations
              </Text>
              <Text
                color="#768BAD"
                fontWeight={400}
                fontSize={12}
                lineHeight="19px"
              >
                <P>
                  2.1 | By agreeing to be bound by these Terms and Conditions,
                  you also agree to be bound by the Duel Rules and Privacy
                  Policy that are hereby incorporated by reference into these
                  Terms and Conditions. In the event of any inconsistency, these
                  Terms and Conditions will prevail. You hereby represent and
                  warrant that:
                </P>
                <P>
                  2.1.1 | You are over (a) 18 and (b) such other legal age or
                  age of majority as determined by any laws which are applicable
                  to you, whichever age is greater;
                </P>
                <P>
                  2.1.2 | You have full capacity to enter into a legally binding
                  agreement with us and you are not restricted by any form of
                  limited legal capacity;
                </P>
                <P>
                  2.1.3 | You participate in the Games strictly in your personal
                  and non-professional capacity; and participate for
                  recreational and entertainment purposes only;
                </P>
                <P>
                  2.1.4 | You participate in the Games on your own behalf and
                  not on the behalf of any other person;
                </P>
                <P>
                  2.1.5 | All information that you provide to us during the term
                  of validity of this agreement is true, complete, correct, and
                  that you shall immediately notify us of any change of such
                  information;
                </P>
                <P>
                  2.1.6 | You are solely responsible for reporting and
                  accounting for any taxes applicable to you under relevant laws
                  for any winnings that you receive from us;
                </P>
                <P>
                  2.1.7 | You understand that by using our services you take the
                  risk of losing money deposited into your Member Account and
                  accept that you are fully and solely responsible for any such
                  loss;
                </P>
                <P>
                  2.1.8 | You are permitted in the jurisdiction in which you are
                  located to use online casino services;
                </P>
                <P>
                  2.1.9 | You will not use our services while located in any
                  jurisdiction that prohibits the placing and/or accepting of
                  bets online (incl. denominated in Bitcoin or any other
                  cryptocurrencies that we use), and/or playing casino and/or
                  live games including for and/or with Crypto;
                </P>
                <P>
                  2.1.10 | In relation to deposits and withdrawals of funds into
                  and from your Member Account, you shall only use Cryptos that
                  are valid and lawfully belong to you;
                </P>
                <P>
                  2.1.11 | You understand that the value of Cryptocurrencies can
                  change dramatically depending on the market value;
                </P>
                <P>
                  2.1.12 | The computer software, the computer graphics, the
                  Websites and the user interface that we make available to you
                  is owned by Duel or its associates and is protected by
                  copyright laws. You may only use the software for your own
                  personal, recreational uses in accordance with all rules,
                  terms and conditions we have established and in accordance
                  with all applicable laws, rules and regulations;
                </P>
                <P>
                  2.1.13 | You understand that Crypto is not considered a legal
                  currency or tender and as such on the Website they are treated
                  as virtual funds with no intrinsic value.
                </P>
                <P>
                  2.1.14 | You affirm that you are not an officer, director,
                  employee or working for any company related to Duel, or a
                  relative or spouse of any of the foregoing;
                </P>
                <P>
                  2.1.15 | You are not diagnosed or classified as a compulsive
                  or problem gambler. We are not accountable if such problem
                  gambling arises whilst using our services, but will endeavor
                  to inform of relevant assistance available. We reserve the
                  right to implement cool off periods if we believe such actions
                  will be of benefit.
                </P>
                <P>
                  2.1.16 | You are not politically exposed person or a family
                  member of a politically exposed person;
                </P>
                <P>
                  2.1.17 | You have only one Member Account with us and agree to
                  not to open any more Member Accounts with us;
                </P>
                <P>
                  2.1.18 | You accept and acknowledge that we reserve the right
                  to detect and prevent the use of prohibited techniques,
                  including but not limited to fraudulent transaction detection,
                  automated registration and signup, gameplay and screen capture
                  techniques. These steps may include, but are not limited to,
                  examination of Players device properties, detection of
                  geo-location and IP masking, transactions and blockchain
                  analysis;
                </P>
                <P>
                  2.1.19 | You accept our right to terminate and/or change any
                  games or events being offered on the Website, and to refuse
                  and/or limit bets.
                </P>
              </Text>
            </Box>

            <Box my={15}>
              <Text color="#768BAD" fontWeight={600} fontSize={16}>
                3. Your Member Account
              </Text>
              <Text
                color="#768BAD"
                fontWeight={400}
                fontSize={12}
                lineHeight="19px"
              >
                <P>
                  3.1 | In order for you to be able to place bets on our
                  websites, you must first personally register a Member Account
                  with us.
                </P>
                <P>
                  3.2 | We do not wish to and shall not accept registration from
                  persons resident in jurisdictions that prohibit you from
                  participating in online gambling, gaming, and/or games of
                  skill, for and/or with cryptocurrencies. By registering a
                  Member Account with us you confirm that you&apos;re not using
                  any third party software to access our sites from
                  jurisdictions that are prohibited, please refer to section 3.3
                  for the jurisdictions that are prohibited.
                </P>
                <P>
                  3.3 | You are aware that the right to access and use the
                  website and any products there offered, may be considered
                  illegal in certain countries. We are not able to verify the
                  legality of service in each and every jurisdiction,
                  consequently, you are responsible in determining whether your
                  accessing and using our website is compliant with the
                  applicable laws in your country and you warrant to us that
                  gambling is not illegal in the territory where you reside. For
                  various legal or commercial reasons, we do not permit Member
                  Accounts to be opened or used by customers resident in certain
                  jurisdictions, including Afghanistan, Australia, Belarus,
                  Belgium, Côte d&apos;Ivoire, Cuba, Curaçao, Czech Republic,
                  Democratic Republic of the Congo, France, Germany, Greece,
                  Iran, Iraq, Italy, Liberia, Libya, Lithuania, Netherlands,
                  North Korea, Portugal, Serbia, Slovakia, South Sudan, Spain,
                  Sudan, Sweden, Syria, United Kingdom, United States of
                  America, Zimbabwe (the &quot;Prohibited Jurisdictions&quot;)
                  are not permitted make use of the Service. By using the
                  Website you confirm you are not a resident in a Restricted
                  Jurisdiction.
                </P>
                <P>
                  3.4 | When attempting to open a Member Account or using the
                  Website, it is the responsibility of the player to verify
                  whether gambling is legal in that particular jurisdiction. If
                  you open or use the Website while residing in a Restricted
                  Jurisdiction: your Member Account may be closed by us
                  immediately; any winnings and rewards will be confiscated and
                  remaining balance returned (subject to reasonable charges),
                  and any returns, winnings or rewards which you have gained or
                  accrued will be forfeited by you and may be reclaimed by us;
                  and you will return to us on demand any such funds which have
                  been withdrawn.
                </P>
                <P>
                  3.5 | You are allowed to have only one Member Account. If you
                  attempt to open more than one Member Account, all of your
                  Member Accounts may be blocked, suspended or closed and any
                  cryptocurrencies credited to your Member Account(s) will be
                  frozen.
                </P>
                <P>
                  3.6 | If you notice that you have more than one registered
                  Member Account you must notify us immediately. Failure to do
                  so may lead to your Member Account being blocked.
                </P>
                <P>
                  3.7 | You will inform us as soon as you become aware of any
                  errors with respect to your Member Account or any calculations
                  with respect to any bet you have placed. We reserve the right
                  to declare null and void any bets that are subject to such an
                  error.
                </P>
                <P>
                  3.8 | If you do not use your Member Account for a time period
                  of 6 months, you will receive a notice from us. If your Member
                  Account remains dormant and unused after this notice for a
                  period of 6 months, we reserve the right to deduct monthly
                  administrative costs from the remaining balance in your Member
                  Account up to a maximum value of 2.5% per month of inactivity
                  from any funds that are remaining in your Member Account to
                  increase security on funds. If this happens, contact us at
                  support@duelana.com to reopen your Member Account.
                </P>
                <P>
                  3.9 | You must enter all mandatory information requested into
                  the registration form, including a valid email address. If you
                  do not enter a valid email address, we will be unable to help
                  you recover any “forgotten passwords”. It is your sole
                  responsibility to ensure that the information you provide is
                  true, complete and correct.
                </P>
                <P>
                  3.10 | We have the right to carry out “KYC” (Know Your
                  Customer) verification procedures and access to your Member
                  Account may be blocked or closed if we determine that you have
                  supplied false or misleading information.
                </P>
                <P>
                  3.11 | As part of the registration process, you will have to
                  choose a username and password for your login into the
                  Website(s). You will have to choose a username which is not
                  disruptive or offensive. It is your sole and exclusive
                  responsibility to ensure that your login details are kept
                  securely. You must not disclose your login details to anyone.
                  We are not liable or responsible for any abuse or misuse of
                  your Member Account by third parties due to your disclosure,
                  whether intentional, accidental, active or passive, of your
                  login details to any third party.
                </P>
                <P>
                  3.12 | If you change your password, you will be unable to
                  withdraw for 48 hours due to security reasons.
                </P>
              </Text>
            </Box>

            <Box>
              <Text color="#768BAD" fontWeight={600} fontSize={16}>
                4. Deposits
              </Text>
              <Text
                color="#768BAD"
                fontWeight={400}
                fontSize={12}
                lineHeight="19px"
              >
                <P>
                  4.1 | You may participate in any Game only if you have
                  sufficient currency balance on your Member Account for such
                  participation. We shall not give you any credit whatsoever for
                  participation in any Game.
                </P>
                <P>
                  4.2 | To deposit funds into your Member Account, you can
                  transfer funds from crypto-wallets and credit cards under your
                  control. Deposits can only be made with your own funds.
                </P>
                <P>
                  4.3 | We reserve the right to use additional procedures and
                  means to verify your identity when processing deposits into
                  your Member Account.
                </P>
                <P>
                  4.4 | Note that some payment methods may include an additional
                  fee. In this case, the fee will be clearly visible for you in
                  the cashier.
                </P>
                <P>
                  4.5 | Note that your bank or payment service provider may
                  charge you additional fees for deposits, withdrawals of
                  currency conversion according to their terms and conditions
                  and your user agreement.
                </P>
              </Text>
            </Box>

            <Box my={15}>
              <Text color="#768BAD" fontWeight={600} fontSize={16}>
                5. Withdrawals
              </Text>
              <Text
                color="#768BAD"
                fontWeight={400}
                fontSize={12}
                lineHeight="19px"
              >
                <P>
                  5.1 | All withdrawals shall be processed in accordance with
                  our withdrawal policy. Crypto withdrawals will be made to your
                  stated Crypto wallet address when making a valid withdrawal
                  request. To withdraw any funds which have been deposited, we
                  require there to be at least 3 blockchain confirmations of the
                  deposit before a withdrawal can be requested.
                </P>
                <P>
                  5.2 | If we mistakenly credit your Member Account with
                  winnings that do not belong to you, whether due to a technical
                  error in the pay-tables, or human error or otherwise, the
                  amount will remain our property and will be deducted from your
                  Member Account. If you have withdrawn funds that do not belong
                  to you prior to us becoming aware of the error, the mistakenly
                  paid amount will (without prejudice to other remedies and
                  actions that may be available at law) constitute a debt owed
                  by you to us. In the event of an incorrect crediting, you are
                  obliged to notify us immediately by email.
                </P>
                <P>
                  5.3 | Duel reserves the right to carry out additional KYC
                  verification procedures for any withdrawals exceeding the
                  equivalent of 1 Bitcoin or $2000 as regulated by our gaming
                  license, and further reserves the right to carry out such
                  verification procedures in case of smaller withdrawals, as
                  demanded by our gaming license. Member Account Holders who
                  wish to recover funds held in a closed, locked or excluded
                  Member Account, are advised to contact Customer Support.
                </P>
                <P>
                  5.4 | All transactions shall be checked in order to prevent
                  money laundering. If the Member becomes aware of any
                  suspicious activity relating to any of the Games of the
                  Website, s/he must report this to Duel immediately. Duel may
                  suspend, block or close a Member Account and withhold funds if
                  requested to do so in accordance with the Prevention of Money
                  Laundering Act or on any other legal basis requested by any
                  state authority. Enhanced due diligence may be done in respect
                  of withdrawals of funds not used for wagering.
                </P>
                <P>
                  5.5 | We reserve the right to apply a wagering requirement of
                  at least 5 (five) times the deposit amount if we suspect the
                  player is using our service as a mixer. It is strictly
                  forbidden to use our service for any other purpose than
                  entertainment.
                </P>
                <P>
                  5.6 | You acknowledge that the funds in your Member Account
                  are consumed instantly when playing and we do not provide
                  return of goods, refunds or retrospective cancellation of your
                  Member Account.
                </P>
                <P>
                  5.7 | If you win 25 bitcoins or more, we reserve the right to
                  pay a maximum of up to 25 bitcoins per week until the full
                  amount is settled.
                </P>
                <P>
                  5.8 | You will not earn any interest on outstanding amounts
                  and acknowledge that the Company is not a financial
                  institution.
                </P>
                <P>
                  5.9 | If you are eligible for a reward, for example a login
                  reward or a deposit reward of 100% up to a certain amount,
                  wagering requirements will apply before you are eligible to
                  make any cash-outs of the reward or winnings. The wagering
                  requirements, which can vary, will be displayed when receiving
                  the reward. If you would like to request a withdrawal before
                  the wagering requirements are fulfilled, Duel will deduct the
                  whole reward amount as well as any winnings before approving
                  the withdrawal. Duel reserves the right to impose, at our own
                  discretion, geographical limitations to individual reward
                  schemes. Local wagering requirements may be applied.
                  rewards/free spins at Duel can only be received once per
                  household/IP (several Member Accounts registered with the same
                  IP address). Duel wagering requirements do not apply to risk
                  free bets.
                </P>
                <P>
                  5.10 | You must use your reward and/ or reward program within
                  30 days from receiving the reward on your Member Account. When
                  the reward and/ or reward program has not been used within 30
                  days from receiving it, Duel reserves the right to cancel any
                  such reward and/ or reward program and may deduct the reward
                  or reward-like reward or freespin immediately after the lapse
                  of the 30 day period.
                </P>
                <P>
                  5.11 | You acknowledge and understand that separate terms and
                  conditions exist with respect to promotions, rewards and
                  special offers, and are in addition to these terms and
                  conditions. These terms and conditions are set forth in the
                  respective content page on this website (https://duel.win), or
                  have been delivered to you personally, as the case may be. In
                  the event of a conflict between the provisions of such
                  promotions, rewards and special offers, and the provisions of
                  these terms and conditions, the provisions of such promotions,
                  rewards and special offers will prevail.
                </P>
                <P>
                  5.12 | We reserve the right to insist that players bet the
                  full amount of their own deposit before they can bet with the
                  free money we credit to them.
                </P>
                <P>
                  5.13 | Certain promotions may be subject to withdrawal and/or
                  cancellation and may only be available for specific periods
                  and on certain specific terms. You must ensure that the
                  promotion you are interested in is still available, that you
                  are eligible, and that you understand any terms which apply to
                  it.
                </P>
                <P>
                  5.14 | Where any term of the offer or promotion is breached or
                  there is any evidence of a series of bets placed by a customer
                  or group of customers, which due to a deposit reward, enhanced
                  payments, free bets, risk free bets or any other promotional
                  offer results in guaranteed customer profits irrespective of
                  the outcome, whether individually or as part of a group, Duel
                  reserves the right to reclaim the reward element of such
                  offers and in their absolute discretion either settle bets at
                  the correct odds, void the free bet reward and risk free bets
                  or void any bet funded by the deposit reward. In addition,
                  Duel reserves the right to levy an administration charge on
                  the customer up to the value of the deposit reward, free bet
                  reward, risk free bet or additional payment to cover
                  administrative costs. We further reserve the right to ask any
                  player to provide sufficient documentation for us to be
                  satisfied in our absolute discretion as to the player’s
                  identity prior to us crediting any reward, free bet, risk free
                  bet or offer to their account.
                </P>
                <P>
                  5.15 | All Duel offers are intended for recreational players
                  and Duel may in its sole discretion limit the eligibility of
                  players to participate in all or part of any promotion.
                </P>
                <P>
                  5.16 | Reward rounds and free spins do not qualify for the
                  jackpot rewards pursuant to casino software provider rules.
                  Only real money rounds qualify for the jackpot rewards.
                </P>
                <P>
                  5.17 | If we determine, in our sole discretion, that you are
                  using the “Double Spend“ methodology, Duel shall void all bets
                  and winnings. Specifically, if you win, then confirm your
                  deposit on the Blockchain and attempt to withdraw, all
                  winnings will be confiscated and your account will be closed
                  permanently. We shall also exercise this right where similar
                  activities are attempted from any connected accounts.
                </P>
                <P>
                  5.18 | Duel reserves the right to amend, cancel, reclaim or
                  refuse any promotion at its own discretion.
                </P>
                <P>
                  5.19 | Note that some payment methods may include an
                  additional fee. In this case, the fee will be clearly visible
                  for you in the cashier.
                </P>
                <P>
                  5.20 | Note that your bank or payment service provider may
                  charge you additional fees for deposits, withdrawals of
                  currency conversion according to their terms and conditions
                  and your user agreement.
                </P>
              </Text>
            </Box>

            <Box>
              <Text color="#768BAD" fontWeight={600} fontSize={16}>
                6. Closing Of Member Accounts
              </Text>
              <Text
                color="#768BAD"
                fontWeight={400}
                fontSize={12}
                lineHeight="19px"
              >
                <P>
                  6.1 | If you wish to close your Member Account, you may do so
                  at any time, by contacting customer support in written form.
                  The effective closure of the Member Account will correspond to
                  the termination of the Terms and Conditions. If the reason
                  behind the closure of the Member Account is related to
                  concerns about possible gambling addiction, you shall indicate
                  this in writing when requesting the closure.
                </P>
              </Text>
            </Box>

            <Box my={15}>
              <Text color="#768BAD" fontWeight={600} fontSize={16}>
                7. Privacy Policy
              </Text>
              <Text
                color="#768BAD"
                fontWeight={400}
                fontSize={12}
                lineHeight="19px"
              >
                <P>
                  7.1 | You hereby acknowledge and accept that if we deem
                  necessary, we are able to collect and otherwise use your
                  personal data in order to allow you access and use of the
                  Websites and in order to allow you to participate in the
                  Games.
                </P>
                <P>
                  7.2 | We hereby acknowledge that in collecting your personal
                  details as stated in the previous provision, we are bound by
                  applicable privacy laws under the laws of Curacao We will
                  protect your personal information and respect your privacy in
                  accordance with best business practices and applicable laws.
                  Please refer to the following link for a fulsome a copy of our
                  privacy policy:{" "}
                  <StyledLink
                    href="https://docs.duel.win/more-info/privacy-policy"
                    rel="noreferrer"
                    target={"_blank"}
                  >
                    https://docs.duel.win/more-info/privacy-policy
                  </StyledLink>
                </P>
                <P>
                  7.3 | We will use your personal data to allow you to
                  participate in the Games and to carry out operations relevant
                  to your participation in the Games. We may also use your
                  personal data to inform you of changes, new services and
                  promotions that we think you may find interesting. If you do
                  not wish to receive such direct marketing correspondences, you
                  may opt out of the service.
                </P>
                <P>
                  7.4 | Your personal data will not be disclosed to third
                  parties, unless such disclosure is necessary for the
                  processing of your requests in relation to your participation
                  in the Games or unless it is required by law. As Duel&apos;s
                  business partners or suppliers or service providers may be
                  responsible for certain parts of the overall functioning or
                  operation of the Website, personal data may be disclosed to
                  them. The employees of Duel have access to your personal data
                  for the purpose of executing their duties and providing you
                  with the best possible assistance and service. You hereby
                  consent to such disclosures.
                </P>
                <P>
                  7.5 | We shall keep all information provided as personal data.
                  You have the right to access personal data held by us about
                  you. No data shall be destroyed unless required by law, or
                  unless the information held is no longer required to be kept
                  for the purpose of the relationship.
                </P>
                <P>
                  7.6 | In order to make your visit to the Websites more
                  user-friendly, to keep track of visits to the Websites and to
                  improve the service, we collect a small piece of information
                  sent from your browser, called a cookie. You can, if you wish,
                  turn off the collection of cookies. You must note, however,
                  that turning off cookies may severely restrict or completely
                  hinder your use of the Websites.
                </P>
              </Text>
            </Box>

            <Box>
              <Text color="#768BAD" fontWeight={600} fontSize={16}>
                8. General Betting Rules
              </Text>
              <Text
                color="#768BAD"
                fontWeight={400}
                fontSize={12}
                lineHeight="19px"
              >
                <P>
                  8.1 | A bet can only be placed by a registered Member Account
                  holder.
                </P>
                <P>8.2 | A bet can only be placed over the internet.</P>
                <P>
                  8.3 | You can only place a bet if you have sufficient balance
                  in your Member Account with Duel.
                </P>
                <P>
                  8.4 | The bet, once concluded, will be governed by the version
                  of the Terms and Conditions valid and available on the Website
                  at the time of the bet being accepted.
                </P>
                <P>
                  8.5 | Any payout of a winning bet is credited to your Member
                  Account, consisting of the stake multiplied by the odds at
                  which the bet was placed.
                </P>
                <P>
                  8.6 | Duel reserves the right to adjust a bet payout credited
                  to a Duel Member Account if it is determined by Duel in its
                  sole discretion that such a payout has been credited due to an
                  error.
                </P>
                <P>
                  8.7 | A bet, which has been placed and accepted, cannot be
                  amended, withdrawn or cancelled by you.
                </P>
                <P>
                  8.8 | The list of all the bets, their status and details are
                  available to you on the Website.
                </P>
                <P>
                  8.9 | When you place a bet you acknowledge that you have read
                  and understood in full all of these Terms and Conditions
                  regarding the bet as stated on the Website.
                </P>
                <P>
                  8.10 | Duel manages your Member Account, calculates the
                  available funds, the pending funds, the betting funds as well
                  as the amount of winnings. Unless proven otherwise, these
                  amounts are considered as final and are deemed to be accurate.
                </P>
                <P>8.11 | You are fully responsible for the bets placed.</P>
                <P>
                  8.12 | Winnings will be paid into your Member Account after
                  the final result is confirmed.
                </P>
              </Text>
            </Box>

            <Box my={15}>
              <Text color="#768BAD" fontWeight={600} fontSize={16}>
                9. Miscarried and Aborted Games
              </Text>
              <Text
                color="#768BAD"
                fontWeight={400}
                fontSize={12}
                lineHeight="19px"
              >
                <P>
                  9.1 | Duel is not liable for any downtime, server disruptions,
                  lagging, or any technical or political disturbance to the game
                  play. Refunds may be given solely at the discretion of
                  Duel&apos;s management.
                </P>
                <P>
                  9.2 | Duel shall accept no liability for any damages or losses
                  which are deemed or alleged to have arisen out of or in
                  connection with the website or its content; including without
                  limitation, delays or interruptions in operation or
                  transmission, loss or corruption of data, communication or
                  lines failure, any person&apos;s misuse of the site or its
                  content or any errors or omissions in content.
                </P>
                <P>
                  9.3 | In the event of a Website malfunction all wagers are
                  void.
                </P>
              </Text>
            </Box>

            <Box>
              <Text color="#768BAD" fontWeight={600} fontSize={16}>
                10. Rewards and Promotions
              </Text>
              <Text
                color="#768BAD"
                fontWeight={400}
                fontSize={12}
                lineHeight="19px"
              >
                <P>
                  10.1 | If you use a deposit reward, no withdrawal of your
                  original deposit will be accepted before you have reached the
                  requirements stipulated under the terms and conditions of the
                  deposit reward.
                </P>
                <P>
                  10.2 | Where any term of the offer or promotion is breached or
                  there is any evidence of a series of bets placed by a customer
                  or group of customers, which due to a deposit reward, enhanced
                  payments, free bets, risk free bets or any other promotional
                  offer results in guaranteed customer profits irrespective of
                  the outcome, whether individually or as part of a group, Duel
                  reserves the right to reclaim the reward element of such
                  offers and in their absolute discretion either settle bets at
                  the correct odds, void the free bet reward and risk free bets
                  or void any bet funded by the deposit reward. In addition,
                  Duel reserves the right to levy an administration charge on
                  the customer up to the value of the deposit reward, free bet
                  reward, risk free bet or additional payment to cover
                  administrative costs. We further reserve the right to ask any
                  customer to provide sufficient documentation for us to be
                  satisfied in our absolute discretion as to the customer&apos;s
                  identity prior to us crediting any reward, free bet, risk free
                  bet or offer to their account.
                </P>
                <P>
                  10.3 | All Duel offers are intended for recreational players
                  and Duel may in its sole discretion limit the eligibility of
                  customers to participate in all or part of any promotion.
                </P>
                <P>
                  10.4 | Duel reserves the right to amend, cancel, reclaim or
                  refuse any promotion at its own discretion.
                </P>
                <P>
                  10.5 | You acknowledge and understand that separate terms and
                  conditions exist with respect to promotions, rewards and
                  special offers, and are in addition to these terms and
                  conditions. These Terms and Conditions are set forth in the
                  respective content page on this website, or have been made
                  available to you personally, as the case may be. In the event
                  of a conflict between the provisions of such promotions,
                  rewards and special offers, and the provisions of these terms
                  and conditions, the provisions of such promotions, rewards and
                  special offers will prevail.
                </P>
                <P>
                  10.6 | We may insist that you bet a certain amount of your own
                  deposit before you can bet with any free/reward funds we
                  credit to your Member Account.
                </P>
                <P>
                  10.7 | You accept that certain promotions may be subject to
                  withdrawal restrictions and/or requirements which need to be
                  met before funds credited under the promotion can be
                  withdrawn. Such terms shall be duly published and made
                  available as part of the promotion. If you opt to make a
                  withdrawal before the applicable wagering requirements are
                  fulfilled, we will deduct the whole reward amount as well as
                  any winnings connected with the use of the reward amounts
                  before approving any withdrawal.
                </P>
              </Text>
            </Box>

            <Box my={15}>
              <Text color="#768BAD" fontWeight={600} fontSize={16}>
                11. Live Chat
              </Text>
              <Text
                color="#768BAD"
                fontWeight={400}
                fontSize={12}
                lineHeight="19px"
              >
                <P>
                  11.1 | As part of your use of the Website we may provide you
                  with a live chat facility, which is moderated by us and
                  subject to controls. We reserve the right to review the chat
                  and to keep a record of all statements made on the facility.
                  Your use of the chat facility should be for recreational and
                  socialising purposes. We reserve the right to remove and ban
                  you from the live chat and suspend, block or cancel your
                  Member Account if you:
                </P>
                <P>
                  11.1.1 | make any statements that are sexually explicit or
                  grossly offensive, including expressions of bigotry, racism,
                  hatred or profanity;
                </P>
                <P>
                  11.1.2 | make statements that are abusive, defamatory or
                  harassing or insulting;
                </P>
                <P>
                  11.1.13 | use the chat facility to advertise, promote or
                  otherwise relate to any other online entities;
                </P>
                <P>
                  11.1.14 | make statements about Duel , or any other Internet
                  site(s) connected to the Website that are untrue and/or
                  malicious and/or damaging to Duel;
                </P>
                <P>
                  11.1.15 | use the chat facility to collude, engage in unlawful
                  conduct or encourage conduct we deem seriously inappropriate.
                  Any suspicious chats will be reported to the competent
                  authority.
                </P>
                <P>
                  11.2 | Live Chat is used as a form of communication between us
                  and you and should not be copied or shared with any forums or
                  third parties.
                </P>
              </Text>
            </Box>

            <Box>
              <Text color="#768BAD" fontWeight={600} fontSize={16}>
                12. Intellectual Property
              </Text>
              <Text
                color="#768BAD"
                fontWeight={400}
                fontSize={12}
                lineHeight="19px"
              >
                <P>
                  12.1 | Duel and its licensors are the sole holders of all
                  rights in and to the Website and code, structure and
                  organization, including copyright, trade secrets, intellectual
                  property and other rights. You may not, within the limits
                  prescribed by applicable laws: (a) copy, distribute, publish,
                  reverse engineer, decompile, disassemble, modify, or translate
                  the website; or (b) use the Website in a manner prohibited by
                  applicable laws or regulations (each of the above is an
                  &quot;Unauthorized Use&quot;). Duel reserves any and all
                  rights implied or otherwise, which are not expressly granted
                  to the Player hereunder and retain all rights, title and
                  interest in and to the Website. You agree that you will be
                  solely liable for any damage, costs or expenses arising out of
                  or in connection with the commission by you of any
                  Unauthorized Use. You shall notify Duel immediately upon
                  becoming aware of the commission by any person of any
                  Unauthorized Use and shall provide Duel with reasonable
                  assistance with any investigations it conducts in light of the
                  information provided by you in this respect.
                </P>
                <P>
                  12.2 | The term &quot;Duel&quot;, its domain names and any
                  other trade marks, or service marks used by Duel as part of
                  the Website (the &quot;Trade Marks&quot;), are solely owned by
                  Duel. In addition, all content on the Website, including, but
                  not limited to, the images, pictures, graphics, photographs,
                  animations, videos, music, audio and text (the &quot;Website
                  Content&quot;) belongs to Duel and is protected by copyright
                  and/or other intellectual property or other rights. You hereby
                  acknowledge that by using the Website, you obtain no rights in
                  the Website Content and/or the Trade Marks, or any part
                  thereof. Under no circumstances may you use the Website
                  Content and/or the Trade Marks without Duel&apos;s prior
                  written consent. Additionally, you agree not to do anything
                  that will harm or potentially harm the rights, including the
                  intellectual property rights of Duel.
                </P>
              </Text>
            </Box>

            <Box my={15}>
              <Text color="#768BAD" fontWeight={600} fontSize={16}>
                13. Limitation of Liability
              </Text>
              <Text
                color="#768BAD"
                fontWeight={400}
                fontSize={12}
                lineHeight="19px"
              >
                <P>
                  13.1 | You enter the Website and participate in the Games at
                  your own risk. The Websites and the Games are provided without
                  any warranty whatsoever, whether expressed or implied.
                </P>
                <P>
                  13.2 | Without prejudice to the generality of the preceding
                  provision, we, our directors, employees, partners, service
                  providers:
                </P>
                <P>
                  13.2.1 | do not warrant that the software, Games and the
                  Websites are fit for their purpose;
                </P>
                <P>
                  13.2.2 | do not warrant that the software, Games and the
                  Websites are free from errors;
                </P>
                <P>
                  13.2.3 | do not warrant that the software, Games and the
                  Websites will be accessible without interruptions;
                </P>
                <P>
                  13.2.4 | shall not be liable for any loss, costs, expenses or
                  damages, whether direct, indirect, special, consequential,
                  incidental or otherwise, arising in relation to your use of
                  the Websites or your participation in the Games.
                </P>
                <P>
                  13.3 | You understand and acknowledge that, if there is a
                  malfunction in a Game or its interoperability, any bets made
                  during such a malfunction shall be void. Funds obtained from a
                  malfunctioning Game shall be considered void, as well as any
                  subsequent game rounds with said funds, regardless of what
                  Games are played using such funds.
                </P>
                <P>
                  13.4 | You hereby agree to fully indemnify and hold harmless
                  us, our directors, employees, partners, and service providers
                  for any cost, expense, loss, damages, claims and liabilities
                  howsoever caused that may arise in relation to your use of the
                  Website or participation in the Games.
                </P>
                <P>
                  13.5 | To the extent permitted by law, our maximum liability
                  arising out of or in connection with your use of the Websites,
                  regardless of the cause of actions (whether in contract, tort,
                  breach of warranty or otherwise), will not exceed $100.
                </P>
              </Text>
            </Box>

            <Box>
              <Text color="#768BAD" fontWeight={600} fontSize={16}>
                14. Breaches, Penalties, and Termination
              </Text>
              <Text
                color="#768BAD"
                fontWeight={400}
                fontSize={12}
                lineHeight="19px"
              >
                <P>
                  14.1 | If you breach any provision of these Terms and
                  Conditions or we have a reasonable ground to suspect that you
                  have breached them, we reserve the right to not open, to
                  suspend, or to close your Member Account, or withhold payment
                  of your winnings and apply such funds to any damages due by
                  you.
                </P>
              </Text>
            </Box>

            <Box my={15}>
              <Text color="#768BAD" fontWeight={600} fontSize={16}>
                15. Severability
              </Text>
              <Text
                color="#768BAD"
                fontWeight={400}
                fontSize={12}
                lineHeight="19px"
              >
                <P>
                  15.1 | If any provision of these Terms and Conditions is held
                  to be illegal or unenforceable, such provision shall be
                  severed from these Terms and Conditions and all other
                  provisions shall remain in force unaffected by such severance.
                </P>
              </Text>
            </Box>

            <Box>
              <Text color="#768BAD" fontWeight={600} fontSize={16}>
                16. Assignment
              </Text>
              <Text
                color="#768BAD"
                fontWeight={400}
                fontSize={12}
                lineHeight="19px"
              >
                <P>
                  16.1 | We reserve the right to assign or otherwise lawfully
                  transfer this agreement. You shall not assign or otherwise
                  transfer this agreement.
                </P>
              </Text>
            </Box>

            <Box my={15}>
              <Text color="#768BAD" fontWeight={600} fontSize={16}>
                17. Entire Agreement
              </Text>
              <Text
                color="#768BAD"
                fontWeight={400}
                fontSize={12}
                lineHeight="19px"
              >
                <P>
                  17.1 | These Terms and Conditions constitute the entire
                  agreement between you and us with respect to the Websites and,
                  save in the case of fraud, supersede all prior or
                  contemporaneous communications and proposals, whether
                  electronic, oral or written, between you and us with respect
                  to the Websites.
                </P>
              </Text>
            </Box>

            <Box>
              <Text color="#768BAD" fontWeight={600} fontSize={16}>
                18. Duel RESTRICTIONS
              </Text>
              <Text color="#768BAD" fontWeight={600} fontSize={12}>
                PERSONAL USE. The Service is intended solely for the User&apos;s
                personal use. The User is only allowed to wager for his/her
                personal entertainment. Users may not create multiple accounts
                for the purpose of collusion, sports betting and/or abuse of
                service.
              </Text>
              <Text color="#768BAD" fontWeight={600} fontSize={12}>
                JURISDICTIONS. Persons located in or reside in Afghanistan,
                Australia, Belarus, Belgium, Côte d&apos;Ivoire, Cuba, Curaçao,
                Czech Republic, Democratic Republic of the Congo, France,
                Germany, Greece, Iran, Iraq, Italy, Liberia, Libya, Lithuania,
                Netherlands, North Korea, Portugal, Serbia, Slovakia, South
                Sudan, Spain, Sudan, Sweden, Syria, United Kingdom, United
                States, Zimbabwe (the &quot;Prohibited Jurisdictions&quot;) are
                not permitted make use of the Service. For the avoidance of
                doubt, the foregoing restrictions on engaging in real-money play
                from Prohibited Jurisdictions apply equally to residents and
                citizens of other nations while located in a Prohibited
                Jurisdiction. Any attempt to circumvent the restrictions on play
                by any persons located in a Prohibited Jurisdiction or
                Restricted Jurisdiction is a breach of this Agreement. An
                attempt at circumvention includes, but is not limited to,
                manipulating the information used by Duel to identify your
                location and providing Duel with false or misleading information
                regarding your location or place of residence.
              </Text>
              <Text color="#768BAD" fontWeight={600} fontSize={12}>
                The attempt to manipulate your real location through the use of
                VPN, proxy, or similar services or through the provision of
                incorrect or misleading information about your place of
                residence, with the intent to circumvent geo-blocking or
                jurisdiction restrictions, constitutes a breach of Clause 5 of
                this Terms of Service.
              </Text>
              <Text
                color="#768BAD"
                fontWeight={400}
                fontSize={12}
                lineHeight="19px"
              >
                <P>
                  18.1 | Duel will not permit the Games to be supplied to any
                  entity that operates in any of the below jurisdictions
                  (irrespective of whether or not Duel Games are being supplied
                  by the entity in that jurisdiction) without the appropriate
                  licenses:
                </P>
              </Text>
              <Text color="#768BAD" fontWeight={600} fontSize={12}>
                Afghanistan, Australia, Belarus, Belgium, Côte d&apos;Ivoire,
                Cuba, Curaçao, Czech Republic, Democratic Republic of the Congo,
                France, Germany, Greece, Iran, Iraq, Italy, Liberia, Libya,
                Lithuania, Netherlands, North Korea, Portugal, Serbia, Slovakia,
                South Sudan, Spain, Sudan, Sweden, Syria, United Kingdom, United
                States, Zimbabwe
              </Text>
            </Box>

            <Box my={15}>
              <Text color="#768BAD" fontWeight={600} fontSize={16}>
                19. Applicable Law and Jurisdiction
              </Text>
              <Text
                color="#768BAD"
                fontWeight={400}
                fontSize={12}
                lineHeight="19px"
              >
                <P>
                  19.1 | The Terms and Conditions and any matters relating
                  hereto or to the Company, shall be governed by, and construed
                  in accordance with, the laws of Curaçao. You irrevocably agree
                  that, subject as provided below, the courts of Curaçao shall
                  have exclusive jurisdiction in relation to any claim, dispute
                  or difference concerning the Company, these Terms and
                  Conditions and any matter arising therefrom and irrevocably
                  waive any right that it may have to object to an action being
                  brought in those courts, or to claim that the action has been
                  brought in an inconvenient forum, or that those courts do not
                  have jurisdiction. Nothing in this clause shall limit the
                  right of Duel to take proceedings against you in any other
                  court of competent jurisdiction, nor shall the taking of
                  proceedings in any one or more jurisdictions preclude the
                  taking of proceedings in any other jurisdictions, whether
                  concurrently or not, to the extent permitted by the law of
                  such other jurisdiction.
                </P>
              </Text>
            </Box>

            <Box>
              <Text color="#768BAD" fontWeight={600} fontSize={16}>
                20. Complaints
              </Text>
              <Text
                color="#768BAD"
                fontWeight={400}
                fontSize={12}
                lineHeight="19px"
              >
                <P>
                  20.1 | If you have a complaint to make regarding our services,
                  you may contact our customer support via the Website live chat
                  or by email at support@duelana.com. We will endeavor to
                  resolve the matter promptly.
                </P>
              </Text>
            </Box>
          </Box>
          <Flex
            flexDirection="row"
            alignItems="center"
            justifyContent="right"
            gap={18}
          >
            <DeclineBtn onClick={props.onDismiss}>Decline</DeclineBtn>
            <AcceptBtn onClick={handleAccept}>Accept</AcceptBtn>
          </Flex>
        </Container>
      </Box>
    </Modal>
  );
}

const Container = styled(Box)`
  max-height: 80vh;
  display: flex;
  flex-direction: column;
  gap: 36px;
`;

const List = styled.ul`
  padding-left: 23px;
  margin: 2px 0;
  line-height: 19px;
`;

const StyledLink = styled.a`
  text-decoration: underline !important;
`;

const AcceptBtn = styled(Button)`
  background: #1a5032;
  border-radius: 5px;
  border: 0;
  font-family: "Inter";
  font-style: normal;
  font-weight: 600;
  font-size: 14px;
  line-height: 17px;
  /* identical to box height */

  padding: 10px 30px;

  color: #4fff8b;
  cursor: pointer;
`;

const DeclineBtn = styled(Button)`
  background: #242f42;
  border-radius: 5px;
  border: 0;
  font-family: "Inter";
  font-style: normal;
  font-weight: 600;
  font-size: 14px;
  line-height: 17px;
  color: #768bad;
  /* identical to box height */
  padding: 10px 30px;
  cursor: pointer;
`;
