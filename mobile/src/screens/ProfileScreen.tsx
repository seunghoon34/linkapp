import React, { useState, useEffect } from 'react';
import { View, Text, StyleSheet, TouchableOpacity, TextInput, ScrollView } from 'react-native';
import { useSelector, useDispatch } from 'react-redux';
import { RootState } from '../store/store';
import { updateProfile } from '../store/slices/authSlice';

const ProfileScreen: React.FC = () => {
  const dispatch = useDispatch();
  const user = useSelector((state: RootState) => state.auth.user);
  const [isEditing, setIsEditing] = useState(false);
  const [profile, setProfile] = useState({
    username: user?.username || '',
    email: user?.email || '',
    bio: user?.bio || '',
    age: user?.age || '',
    gender: user?.gender || '',
    interests: user?.interests || [],
  });

  const handleSave = () => {
    dispatch(updateProfile(profile));
    setIsEditing(false);
  };

  return (
    <ScrollView style={styles.container}>
      <Text style={styles.title}>Profile</Text>
      {isEditing ? (
        <>
          <TextInput
            style={styles.input}
            value={profile.username}
            onChangeText={(text) => setProfile({ ...profile, username: text })}
            placeholder="Username"
          />
          <TextInput
            style={styles.input}
            value={profile.email}
            onChangeText={(text) => setProfile({ ...profile, email: text })}
            placeholder="Email"
            keyboardType="email-address"
          />
          <TextInput
            style={styles.input}
            value={profile.bio}
            onChangeText={(text) => setProfile({ ...profile, bio: text })}
            placeholder="Bio"
            multiline
          />
          <TextInput
            style={styles.input}
            value={profile.age}
            onChangeText={(text) => setProfile({ ...profile, age: text })}
            placeholder="Age"
            keyboardType="numeric"
          />
          <TextInput
            style={styles.input}
            value={profile.gender}
            onChangeText={(text) => setProfile({ ...profile, gender: text })}
            placeholder="Gender"
          />
          <TextInput
  style={styles.input}
  value={profile.interests ? profile.interests.join(', ') : ''}
  onChangeText={(text) => setProfile({ ...profile, interests: text.split(',').map(i => i.trim()) })}
  placeholder="Interests (comma-separated)"
/>
          <TouchableOpacity style={styles.button} onPress={handleSave}>
            <Text style={styles.buttonText}>Save</Text>
          </TouchableOpacity>
        </>
      ) : (
        <>
          <Text style={styles.info}>Username: {profile.username}</Text>
          <Text style={styles.info}>Email: {profile.email}</Text>
          <Text style={styles.info}>Bio: {profile.bio}</Text>
          <Text style={styles.info}>Age: {profile.age}</Text>
          <Text style={styles.info}>Gender: {profile.gender}</Text>
          <Text style={styles.info}>Interests: {profile.interests.join(', ')}</Text>
          <TouchableOpacity style={styles.button} onPress={() => setIsEditing(true)}>
            <Text style={styles.buttonText}>Edit Profile</Text>
          </TouchableOpacity>
        </>
      )}
    </ScrollView>
  );
};

const styles = StyleSheet.create({
  container: {
    flex: 1,
    padding: 20,
  },
  title: {
    fontSize: 24,
    fontWeight: 'bold',
    marginBottom: 20,
  },
  info: {
    fontSize: 18,
    marginBottom: 10,
  },
  input: {
    borderWidth: 1,
    borderColor: '#ccc',
    borderRadius: 5,
    padding: 10,
    marginBottom: 10,
    fontSize: 16,
  },
  button: {
    backgroundColor: '#007AFF',
    padding: 15,
    borderRadius: 5,
    marginTop: 20,
  },
  buttonText: {
    color: 'white',
    fontSize: 16,
    fontWeight: 'bold',
    textAlign: 'center',
  },
});

export default ProfileScreen;