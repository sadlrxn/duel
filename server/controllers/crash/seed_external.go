package crash

/*
/* @External
/* Admin determines original salt and generate whole seed chain.
*/
func DetermineSaltForSeedChain(salt string, length int) (string, error) {
	return generateSeedChainWithSalt(salt, length)
}

/*
/* @External
/* Admin determines client seed and calculate outcomes for each seed.
*/
func DetermineClientSeed(clientSeed string, houseEdge int64, startIndex int) error {
	return determineClientSeed(clientSeed, houseEdge, startIndex)
}
