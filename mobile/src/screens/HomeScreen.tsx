import React, { useState } from 'react';
import { View, Text, StyleSheet, TouchableOpacity } from 'react-native';
import { NativeStackNavigationProp } from '@react-navigation/native-stack';

type RootStackParamList = {
  Home: undefined;
  Login: undefined;
  Register: undefined;
};

type HomeScreenNavigationProp = NativeStackNavigationProp<RootStackParamList, 'Home'>;

type Props = {
  navigation: HomeScreenNavigationProp;
};

const HomeScreen: React.FC<Props> = ({ navigation }) => {
  const [isSearching, setIsSearching] = useState(false);

  const handleStartMatching = () => {
    setIsSearching(true);
    // TODO: Implement API call to start searching
  };

  return (
    <View style={styles.container}>
      <Text style={styles.title}>Welcome to LinkApp</Text>
      <Text style={styles.subtitle}>Find your perfect match with just a tap!</Text>
      {!isSearching ? (
        <TouchableOpacity 
          style={[styles.button, styles.matchingButton]} 
          onPress={handleStartMatching}
        >
          <Text style={styles.buttonText}>Start Matching</Text>
        </TouchableOpacity>
      ) : (
        <Text style={styles.searchingText}>Searching for matches...</Text>
      )}
      <TouchableOpacity 
        style={styles.button} 
        onPress={() => navigation.navigate('Login')}
      >
        <Text style={styles.buttonText}>Login</Text>
      </TouchableOpacity>
      <TouchableOpacity 
        style={[styles.button, styles.registerButton]} 
        onPress={() => navigation.navigate('Register')}
      >
        <Text style={styles.buttonText}>Register</Text>
      </TouchableOpacity>
    </View>
  );
};

const styles = StyleSheet.create({
  container: {
    flex: 1,
    justifyContent: 'center',
    alignItems: 'center',
    padding: 20,
  },
  title: {
    fontSize: 24,
    fontWeight: 'bold',
    marginBottom: 10,
  },
  subtitle: {
    fontSize: 16,
    marginBottom: 20,
    textAlign: 'center',
  },
  button: {
    backgroundColor: '#007AFF',
    padding: 15,
    borderRadius: 5,
    width: '100%',
    alignItems: 'center',
    marginBottom: 10,
  },
  registerButton: {
    backgroundColor: '#34C759',
  },
  buttonText: {
    color: 'white',
    fontSize: 16,
    fontWeight: 'bold',
  },
  matchingButton: {
    backgroundColor: '#4CAF50',
    marginBottom: 20,
  },
  searchingText: {
    fontSize: 18,
    color: '#4CAF50',
    marginBottom: 20,
  },
});

export default HomeScreen;