import { Ionicons } from "@expo/vector-icons"
import { Image, StyleSheet } from "react-native"
import React from "react";
import { useUser } from "../context/UserContext";

const ProfileIcon = ({ color }: { color: string }) => {
  const { user } = useUser();
  
 
  return user?.profileImg ? (
    <Image source={{ uri: user.profileImg }} style={styles.userImage} />
  ) : (
    <Image source={{ uri: "https://xsgames.co/randomusers/avatar.php?g=male" }} style={styles.userImage} />
  );
};

export const icon = {
    index: ({ color }: { color: string }) => (
        <Ionicons name='home-outline' size={22} color={color} />
    ),
    explore: ({ color }: { color: string }) => (
        <Ionicons name='search-outline' size={22} color={color} />
    ),
    notifications: ({ color }: { color: string }) => (
        <Ionicons name='notifications-outline' size={22} color={color} />
    ),
    cart: ({ color }: { color: string }) => (
        <Ionicons name='cart-outline' size={22} color={color} />
    ),
    profile: ({ color }: { color: string }) => (
        <ProfileIcon color={color} />
    ),
}

const styles = StyleSheet.create({
    userImage: {
        width: 24,
        height: 24,
        borderRadius: 20
    },
});