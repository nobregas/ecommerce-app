import { StyleSheet, Text, View, TouchableOpacity, ScrollView } from 'react-native'
import React from 'react'
import { useHeaderHeight } from '@react-navigation/elements'
import { Link, Stack } from 'expo-router'
import { Colors } from '@/constants/Colors'
import { Ionicons } from '@expo/vector-icons'
import { Image } from "expo-image"


type Props = {}

const ProfileScreen = (props: Props) => {
  const headerHeight = useHeaderHeight()


  return (
    <>
      <Stack.Screen
        options={{ headerShown: true, headerTransparent: true, headerTitleAlign: 'center', title: 'Profile' }}
      />
      <View style={[styles.container, { marginTop: headerHeight }]}>
        <View style={{ alignItems: "center" }}>
          <Image source={{ uri: "https://xsgames.co/randomusers/avatar.php?g=male" }} style={styles.profileImg} />
          <Text style={styles.username}>John Doe</Text>
        </View>

        <ScrollView style={styles.buttonWrapper}>
          <TouchableOpacity style={styles.button}>
            <Ionicons name="person-outline" size={20} color={Colors.black} />
            <Text style={styles.btnText}>Your Orders</Text>
          </TouchableOpacity>
          <TouchableOpacity style={styles.button}>
            <Ionicons name="heart-outline" size={20} color={Colors.black} />
            <Text style={styles.btnText}>Your WithList</Text>
          </TouchableOpacity>
          <TouchableOpacity style={styles.button}>
            <Ionicons name="card-outline" size={20} color={Colors.black} />
            <Text style={styles.btnText}>Payment History</Text>
          </TouchableOpacity>
          <TouchableOpacity style={styles.button}>
            <Ionicons name="help-circle-outline" size={20} color={Colors.black} />
            <Text style={styles.btnText}>Customer Suport</Text>
          </TouchableOpacity>
          <TouchableOpacity style={styles.button}>
            <Ionicons name="pencil-outline" size={20} color={Colors.black} />
            <Text style={styles.btnText}>Edit Profile</Text>
          </TouchableOpacity>
          <TouchableOpacity style={styles.button}>
            <Ionicons name="settings-outline" size={20} color={Colors.black} />
            <Text style={styles.btnText}>Settings</Text>
          </TouchableOpacity>
          <Link href="/signin" asChild>
            <TouchableOpacity style={styles.button}>
              <Ionicons name="log-out-outline" size={20} color={Colors.black} />
              <Text style={styles.btnText}>Logout</Text>
            </TouchableOpacity>
          </Link>
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
  }
})