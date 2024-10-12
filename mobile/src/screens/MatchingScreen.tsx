import React, { useEffect, useState } from 'react';
import { View, Text, StyleSheet, ActivityIndicator } from 'react-native';
import { useDispatch, useSelector } from 'react-redux';
import { AppDispatch, RootState } from '../store/store';
import { findMatch } from '../store/slices/matchSlice';

const MatchingScreen: React.FC = () => {
  const dispatch = useDispatch<AppDispatch>();
  const [isSearching, setIsSearching] = useState(true);
  const currentMatch = useSelector((state: RootState) => state.match.currentMatch);

  useEffect(() => {
    const searchForMatch = async () => {
      try {
        await dispatch(findMatch()).unwrap();
        setIsSearching(false);
      } catch (error) {
        console.error('Error finding match:', error);
        setIsSearching(false);
      }
    };

    searchForMatch();
  }, [dispatch]);

  if (isSearching) {
    return (
      <View style={styles.container}>
        <ActivityIndicator size="large" color="#0000ff" />
        <Text style={styles.searchingText}>Searching for your perfect match...</Text>
      </View>
    );
  }

  return (
    <View style={styles.container}>
      {currentMatch ? (
        <Text style={styles.matchText}>You've been matched with {currentMatch.name}!</Text>
      ) : (
        <Text style={styles.noMatchText}>No match found. Try again later.</Text>
      )}
    </View>
  );
};

const styles = StyleSheet.create({
  container: {
    flex: 1,
    justifyContent: 'center',
    alignItems: 'center',
  },
  searchingText: {
    marginTop: 20,
    fontSize: 18,
  },
  matchText: {
    fontSize: 20,
    fontWeight: 'bold',
  },
  noMatchText: {
    fontSize: 18,
    color: 'red',
  },
});

export default MatchingScreen;