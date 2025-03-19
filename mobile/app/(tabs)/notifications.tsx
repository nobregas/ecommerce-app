import { FlatList, StyleSheet, Text, View } from 'react-native'
import React, { useEffect, useState } from 'react'
import { useHeaderHeight } from '@react-navigation/elements'
import { Stack } from 'expo-router'
import { NotificationType } from '@/types/type'
import { Ionicons } from '@expo/vector-icons'
import { Colors } from '@/constants/Colors'
import Animated, { FadeInDown } from 'react-native-reanimated'
import { getNotifications } from '@/service/ApiService'

type Props = {}

const NotificationsScreen = (props: Props) => {
  const headerHeight = useHeaderHeight()

  const [notifications, setNotifications] = useState<NotificationType[]>([])

  useEffect(() => {
    fetchNotifications()
  }, [])

  const fetchNotifications = async () => {
      try {
        const categoriesData = await getNotifications()
  
        setNotifications(categoriesData)
  
      } catch (error) {
        console.log(error)
      }
    }
  

  return (
    <>
      <Stack.Screen
        options={{ headerShown: true, headerTransparent: true, headerTitleAlign: 'center', title: 'Notifications' }}
      />
      <View style={[styles.container, { marginTop: headerHeight }]}>
        <FlatList
          data={notifications}
          keyExtractor={(item) => item.id.toString()}
          renderItem={({ item, index }) => (
            <Animated.View 
            style={styles.notificationWrapper} 
            entering={FadeInDown.delay(300 + index * 100).duration(500)}
            >
              <View style={styles.notificationIcon}>
                <Ionicons name="notifications-outline" size={20} color={Colors.black} />
              </View>
              <View style={styles.notificationInfo}>
                <View style={{ flexDirection: 'row', justifyContent: 'space-between', alignItems: 'center' }}>
                  <Text style={styles.notificationTitle}>{item.title}</Text>
                  <Text style={styles.notificationMessage}>{item.timestamp}</Text>
                </View>
                <Text style={styles.notificationMessage}>{item.message}</Text>
              </View>
            </Animated.View>
          )}
        />
      </View>
    </>
  )
}

export default NotificationsScreen

const styles = StyleSheet.create({
  container: {
    flex: 1,
    paddingHorizontal: 20
  },
  notificationWrapper: {
    flexDirection: 'row',
    alignItems: 'center',
    padding: 10,
    marginBottom: 10,
    borderWidth: StyleSheet.hairlineWidth,
    borderColor: Colors.lightGray,
    backgroundColor: Colors.extraLightGray,
    borderRadius: 5
  },
  notificationIcon: {
    width: 40,
    height: 40,
    justifyContent: 'center',
    alignItems: 'center',
  },
  notificationInfo: {
    flex: 1,
  },
  notificationTitle: {
    fontSize: 16,
    fontWeight: '500',
    color: Colors.black
  },
  notificationMessage: {
    fontSize: 14,
    color: Colors.gray,
    marginTop: 5,
    lineHeight: 20
  },
  
})