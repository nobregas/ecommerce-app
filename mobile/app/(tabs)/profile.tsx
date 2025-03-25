import { StyleSheet, Text, View, TouchableOpacity, ScrollView, ActivityIndicator } from 'react-native'
import React, { useEffect } from 'react'
import { useHeaderHeight } from '@react-navigation/elements'
import { Stack, useRouter } from 'expo-router'
import { Colors } from '@/constants/Colors'
import { Ionicons } from '@expo/vector-icons'
import { Image } from "expo-image"
import { useUser } from '@/context/UserContext'

type Props = {}

const ProfileScreen = (props: Props) => {
  const headerHeight = useHeaderHeight()
  const { user, isLoading, error, fetchUser, logout } = useUser();
  const router = useRouter();

  useEffect(() => {
    if (!user && !isLoading) {
      fetchUser();
    }
  }, [user, isLoading]);

  const handleLogout = async () => {
    await logout();
    router.replace('/(auth)/signin');
  }

  const handleRetry = () => {
    fetchUser();
  };

  if (isLoading) {
    return (
      <View style={styles.loadingContainer}>
        <ActivityIndicator size="large" color={Colors.primary} />
      </View>
    );
  }

  if (error) {
    return (
      <View style={styles.errorContainer}>
        <Text style={styles.errorText}>{error}</Text>
        <TouchableOpacity style={styles.retryButton} onPress={handleRetry}>
          <Text style={styles.retryText}>Try again</Text>
        </TouchableOpacity>
      </View>
    );
  }

  if (!user) {
    return (
      <View style={styles.errorContainer}>
        <Text style={styles.errorText}>Isn't possible to load user data</Text>
        <TouchableOpacity style={styles.retryButton} onPress={handleRetry}>
          <Text style={styles.retryText}>Try again</Text>
        </TouchableOpacity>
      </View>
    );
  }

  return (
    <>
      <Stack.Screen
        options={{ headerShown: true, headerTransparent: true, headerTitleAlign: 'center', title: 'Perfil' }}
      />
      <View style={[styles.container, { marginTop: headerHeight }]}>
        <View style={{ alignItems: "center" }}>
          <Image source={{ uri: user.profileImg || "https://xsgames.co/randomusers/avatar.php?g=male" }} style={styles.profileImg} />
          <Text style={styles.username}>{user.fullName || "Usuário"}</Text>
          <Text style={styles.email}>{user.email}</Text>
        </View>

        <ScrollView style={styles.buttonWrapper}>
          <TouchableOpacity style={styles.button}>
            <Ionicons name="person-outline" size={20} color={Colors.black} />
            <Text style={styles.btnText}>Your Orders</Text>
          </TouchableOpacity>
          <TouchableOpacity style={styles.button}>
            <Ionicons name="heart-outline" size={20} color={Colors.black} />
            <Text style={styles.btnText}>Wishlist</Text>
          </TouchableOpacity>
          <TouchableOpacity style={styles.button}>
            <Ionicons name="card-outline" size={20} color={Colors.black} />
            <Text style={styles.btnText}>Payment History</Text>
          </TouchableOpacity>
          <TouchableOpacity style={styles.button}>
            <Ionicons name="help-circle-outline" size={20} color={Colors.black} />
            <Text style={styles.btnText}>Customer Support</Text>
          </TouchableOpacity>
          <TouchableOpacity style={styles.button}>
            <Ionicons name="pencil-outline" size={20} color={Colors.black} />
            <Text style={styles.btnText}>Edit Profile</Text>
          </TouchableOpacity>
          <TouchableOpacity style={styles.button}>
            <Ionicons name="settings-outline" size={20} color={Colors.black} />
            <Text style={styles.btnText}>Settings</Text>
          </TouchableOpacity>
          <TouchableOpacity style={styles.button} onPress={handleLogout}>
            <Ionicons name="log-out-outline" size={20} color={Colors.black} />
            <Text style={styles.btnText}>Logout</Text>
          </TouchableOpacity>
        </ScrollView>
      </View>
    </>
  )
}

export default ProfileScreen

const styles = StyleSheet.create({
  container: {
    flex: 1,
    padding: 20
  },
  profileImg: {
    width: 100,
    height: 100,
    borderRadius: 50
  },
  username: {
    fontSize: 20,
    fontWeight: "500",
    marginTop: 10,
    color: Colors.black
  },
  buttonWrapper: {
    marginTop: 20,
  },
  button: {
    padding: 10,
    borderColor: Colors.lightGray,
    borderWidth: 1,
    borderRadius: 5,
    flexDirection: "row",
    alignItems: "center",
    gap: 10,
    marginBottom: 10
  },
  btnText: {
    fontSize: 14,
    fontWeight: "500",
    color: Colors.black
  },
  loadingContainer: {
    flex: 1,
    justifyContent: 'center',
    alignItems: 'center'
  },
  errorContainer: {
    flex: 1,
    justifyContent: 'center',
    alignItems: 'center',
    padding: 20
  },
  errorText: {
    fontSize: 16,
    color: 'red',
    textAlign: 'center',
    marginBottom: 15
  },
  retryButton: {
    backgroundColor: Colors.primary,
    padding: 10,
    borderRadius: 5
  },
  retryText: {
    color: 'white',
    fontWeight: '500'
  },
  email: {
    fontSize: 14,
    color: Colors.gray,
    marginTop: 5
  },
})