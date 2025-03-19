import { View, FlatList, StyleSheet, Dimensions, ViewToken } from 'react-native'
import React, { useRef, useState } from 'react'
import Pagination from '@/components/Pagination'
import { Image } from "expo-image"


type Props = {
    images: string[]
}

const width = Dimensions.get('window').width

const ImageSlider = ({ images }: Props) => {

    const [paginationIndex, setPaginationIndex] = useState(0)

    const viewabilityConfig = {
        itemVisiblePercentThreshold: 50,
    }

    const onViewableItemsChanged = ({ viewableItems }: { viewableItems: ViewToken[] }) => {
        if (
            viewableItems[0].index !== undefined &&
            viewableItems[0].index !== null
        ) {
            setPaginationIndex(viewableItems[0].index % images.length)
        }
    }

    const viewabilityConfigCallbackPairs = useRef([
        { viewabilityConfig, onViewableItemsChanged }
    ])

    return (
        <View >
            <FlatList
                data={images}
                horizontal
                showsHorizontalScrollIndicator={false}
                viewabilityConfigCallbackPairs={viewabilityConfigCallbackPairs.current}
                pagingEnabled
                scrollEventThrottle={16}
                renderItem={({ item }) => (
                    <View style={styles.container}>
                        <Image source={{ uri: item }} style={styles.itemImg} />
                    </View>
                )}
            />
            <Pagination items={images} paginationIndex={paginationIndex} />
        </View>
    )
}

export default ImageSlider

const styles = StyleSheet.create({
    container: {
        width: width,
        justifyContent: 'center',
        alignItems: 'center',
    },
    itemImg: {
        width: 300,
        height: 300,
        borderRadius: 10
    }
});